package dao

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/luci/go-render/render"

	"uuid_server/model"
)

var redisCount = int64(0)

type UUIDRedis struct {
}

func (r *UUIDRedis) GetNextInterval(bizCode int64, interval int64) (int64, int64, model.DataSource, error) {
	end := atomic.AddInt64(&redisCount, interval)
	return end - interval, end, model.RedisSource, nil
}

/*
	先使用硬编码模拟，后续需要使用redis list实现
*/

var lock sync.Mutex

var mysqlBounds = map[int64][]*model.Bound{
	0: {
		&model.Bound{
			Start: 1,
			End:   100,
		},
		&model.Bound{
			Start: 100,
			End:   10000,
		},
		&model.Bound{
			Start: 10001,
			End:   100000,
		},
	},
}

var redisBounds = map[int64][]*model.Bound{
	0: {
		&model.Bound{
			Start: 1,
			End:   100,
		},
		&model.Bound{
			Start: 100,
			End:   10000,
		},
		&model.Bound{
			Start: 10001,
			End:   100000,
		},
	},
}

func (r *UUIDRedis) StoreBounds(bounds []*model.Bound, bizCode int64, dataSource model.DataSource) error {
	fmt.Printf("store bounds:%v \n", render.Render(bounds))
	return nil
}

func (r *UUIDRedis) LoadBounds(bizCode int64, dataSource model.DataSource) ([]*model.Bound, error) {
	lock.Lock()
	defer lock.Unlock()
	switch dataSource {
	case model.MysqlCacheSource:
		return mysqlBounds[bizCode], nil
	case model.RedisCacheSource:
		return redisBounds[bizCode], nil
	}
	fmt.Printf("load bounds err,biCode:%v,dataSource:%v \n", bizCode, dataSource.String())
	return nil, fmt.Errorf("invalid dataSource")
}
