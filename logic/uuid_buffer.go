package logic

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"uuid_server/tools"

	"github.com/luci/go-render/render"

	"uuid_server/dao"
	"uuid_server/model"
	"uuid_server/utils"
)

const (
	defaultCap           int64 = 200000
	defaultFactor        int64 = 75
	defaultQPSWindow     int64 = 30
	defaultQPSBuffer     int64 = 1000
	maxCacheMiss         int64 = 2
	defaultMinFetchCount int64 = 10000
	defaultMaxFetchCount int64 = 1000000
)

type Config struct {
	DaoFetcher    dao.DaoFetcher
	Source        model.DataSource
	BizCode       int64
	Cap           *int64
	CacheMiss     *int64
	Factor        *int64
	MinFetchCount *int64
	MaxFetchCount *int64
	QPSBuffer     *int64
	QPSWindow     *int64
}

func NewDefaultMysqlConfig() *Config {
	return &Config{
		DaoFetcher: &dao.UIDCounterMysql{},
		Source:     model.MysqlCacheSource,
	}
}

func NewDefaultRedisConfig() *Config {
	return &Config{
		DaoFetcher: &dao.UIDCounterRedis{},
		Source:     model.RedisCacheSource,
	}
}

type CheckType int8

func (c CheckType) String() string {
	switch c {
	case ReSize:
		return "resize"
	case CacheMiss:
		return "cache_miss"
	default:
		return "unknown_type"
	}
}

const (
	ReSize CheckType = iota
	CacheMiss
)

type UUIDBuffer struct {
	qpsCounter   *tools.QPSCounter
	checkCapChan chan CheckType // 触发调整容量的类型 1. size/cap少于阙值 2.缓存穿透

	linkedList *tools.LinkedList
	putLock    sync.Mutex // 只有一个协程put，不需要加锁
	getLock    sync.Mutex

	daoFetcher dao.DaoFetcher
	source     model.DataSource

	isClose            int64
	size               int64 // uuid的数量
	cap                int64 // 容量
	cacheMiss          int64 // 获取id数量大于缓存数，会穿透缓存直接请求数据库，连续穿透超过一定次数会扩容
	missMaxFetchNumber int64 // max(记录缓存穿透情况下获取的uuid数量)

	factor        int64 // 负载因子，当数量低于一定比例，拉取uuid
	minFetchCount int64 // 每次拉取数据最小数量
	maxFetchCount int64
	qpsBuffer     int64 // pqs * qpsBuffer = 拉取数据量
	bizCode       int64 // 业务类型
}

func NewUUIDBuffer(config *Config) *UUIDBuffer {
	buffer := &UUIDBuffer{
		linkedList:   tools.NewLinkedList(),
		checkCapChan: make(chan CheckType, 1),

		daoFetcher: &dao.UIDCounterMysql{},
		source:     model.MysqlCacheSource,

		cap:           defaultCap,
		factor:        defaultFactor,
		minFetchCount: defaultMinFetchCount,
		maxFetchCount: defaultMaxFetchCount,
		qpsBuffer:     defaultQPSBuffer,
	}

	buffer.daoFetcher = config.DaoFetcher
	buffer.source = config.Source
	buffer.bizCode = config.BizCode

	qpsWindow := defaultQPSWindow
	if config.QPSWindow != nil {
		qpsWindow = *config.QPSWindow
	}
	buffer.qpsCounter = tools.NewQPSCounter(qpsWindow)

	if config.Cap != nil {
		buffer.cap = *config.Cap
	}

	if config.Factor != nil {
		buffer.factor = *config.Factor
	}

	if config.MinFetchCount != nil {
		buffer.minFetchCount = *config.MinFetchCount
	}

	if config.MaxFetchCount != nil {
		buffer.maxFetchCount = *config.MaxFetchCount
	}

	if config.QPSBuffer != nil {
		buffer.qpsBuffer = *config.QPSBuffer
	}
	return buffer
}

func (p *UUIDBuffer) Start() {
	go p.connFetch()
	// 初始化触发拉取数据
	select {
	case p.checkCapChan <- ReSize:
	default:
	}
}

// connFetch调整缓存容量并拉取数据
// 两种情况会调整容量 1.缓存数据量低于阙值，根据qps调整 2.缓存穿透
func (p *UUIDBuffer) connFetch() {
	handlerFunc := func() error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("connFetch panic msg:%v \n", err)
			}
		}()

		if p.isBufferClose() {
			return errors.New("uuid buffer is closed")
		}

		newCap := p.qpsCounter.SumCount() * p.qpsBuffer
		var checkType CheckType
		select {
		// 主要是为了定时判断buffer关闭
		case <-time.After(10 * time.Second):
			return nil
		case checkType = <-p.checkCapChan:
			switch checkType {
			case ReSize:
			case CacheMiss:
				p.cacheMiss++
				if p.cacheMiss <= maxCacheMiss {
					return nil
				}

				newCap = utils.MaxInt64(newCap, p.missMaxFetchNumber)
				p.cacheMiss = 0
			}
		}

		if newCap < p.minFetchCount {
			newCap = p.minFetchCount
		} else if newCap > p.maxFetchCount {
			newCap = p.maxFetchCount
		}

		atomic.StoreInt64(&p.cap, newCap)
		// 判断是否需要拉取数据
		if p.isNeedExtend() {
			start, end, _, err := p.daoFetcher.GetNextInterval(p.bizCode, newCap)
			if err != nil {
				return err
			}

			p.linkedList.Push(&model.Bound{Start: start, End: end})
			atomic.AddInt64(&p.size, newCap)

			fmt.Printf("fetch data,cap:%v,newCap:%v size:%v \n", atomic.LoadInt64(&p.cap), newCap, atomic.LoadInt64(&p.size))
		}

		fmt.Printf("adjust cap,cap:%v,newCap:%v checkType:%v \n", atomic.LoadInt64(&p.cap), newCap, checkType.String())

		return nil
	}

	for {
		if err := handlerFunc(); err != nil {
			return
		}
	}
}

func (p *UUIDBuffer) isNeedExtend() bool {
	sz := atomic.LoadInt64(&p.size)
	cp := atomic.LoadInt64(&p.cap)

	// uuid的数量小于一定比例
	return sz < cp*p.factor/100
}

func (p *UUIDBuffer) isBufferClose() bool {
	return atomic.LoadInt64(&p.isClose) == 1
}

func (p *UUIDBuffer) Stop() {
	atomic.StoreInt64(&p.isClose, 1)
}

func (p *UUIDBuffer) GetUUIDBound(count int64) ([]*model.UUIDBound, error) {
	if p.isBufferClose() {
		return nil, errors.New("uuid buffer is close")
	}

	p.qpsCounter.AddCount(count)

	sz := atomic.LoadInt64(&p.size)
	// 请求数量超过缓存
	if sz < count {
		uuidBounds, err := p.GetBoundFromFetcher(count)
		fmt.Printf("first check get uuids data:%v, err:%v \n", render.Render(uuidBounds), err)
		return uuidBounds, err
	}

	uuidBounds := make([]*model.UUIDBound, 0)
	handlerFunc := func() error {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("GetUUIDBound panic msg:%v \n", err)
			}
		}()

		p.getLock.Lock()
		defer p.getLock.Unlock()

		if p.isBufferClose() {
			return errors.New("uuid buffer is close")
		}

		// double check
		sz := atomic.LoadInt64(&p.size)
		if sz < count {
			var err error
			uuidBounds, err = p.GetBoundFromFetcher(count)

			fmt.Printf("second check get uuids data:%v, err:%v \n", render.Render(uuidBounds), err)

			if err != nil {
				return err
			}
			return nil
		}

		countVal := count
		for countVal > 0 {
			node := p.linkedList.GetFirstNode()
			if node == nil {
				return errors.New("linked node is nil")
			}
			val, _ := node.GetData().(*model.Bound)
			interval := val.End - val.Start
			if count >= interval {
				p.linkedList.Pop()
				uuidBounds = append(uuidBounds, &model.UUIDBound{IntervalBound: val, Source: p.source})
				countVal = countVal - interval
			} else {
				uuidBounds = append(uuidBounds, &model.UUIDBound{IntervalBound: &model.Bound{Start: val.Start, End: val.Start + countVal}, Source: p.source})
				val.Start = val.Start + countVal
				countVal = 0
			}
		}

		atomic.AddInt64(&p.size, -count)
		if p.isNeedExtend() {
			atomic.StoreInt64(&p.missMaxFetchNumber, count)
			select {
			case p.checkCapChan <- ReSize:
			default:
			}
		}

		fmt.Printf("get cache uuids data:%v \n", render.Render(uuidBounds))

		return nil
	}

	if err := handlerFunc(); err != nil {
		return nil, err
	}
	return uuidBounds, nil
}

func (p *UUIDBuffer) GetBoundFromFetcher(count int64) ([]*model.UUIDBound, error) {
	uuidBounds := make([]*model.UUIDBound, 0)
	start, end, source, err := p.daoFetcher.GetNextInterval(p.bizCode, count)
	if err != nil {
		return nil, err
	}

	uuidBounds = append(uuidBounds, &model.UUIDBound{IntervalBound: &model.Bound{Start: start, End: end}, Source: source})

	// 触发容量调整
	select {
	case p.checkCapChan <- CacheMiss:
	default:
	}
	return uuidBounds, nil
}
