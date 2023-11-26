package tools

import (
	"sync"
	"sync/atomic"
	"time"
)

type Counter struct {
	data int64 // 高32位存放秒级别时间戳，可用到2038年，低32位存放统计数
}

func (c *Counter) Serialize() (int64, int64) {
	dt := atomic.LoadInt64(&c.data)
	return dt >> 32, dt & 0xffffffff
}

func (c *Counter) Deserialize(timestamp, total int64) {
	atomic.StoreInt64(&c.data, (timestamp<<32)|total)
}

type QPSCounter struct {
	size     int64
	counters []*Counter
	lock     sync.Mutex
}

func NewQPSCounter(sz int64) *QPSCounter {
	qpsCounter := &QPSCounter{
		counters: make([]*Counter, sz),
		size:     sz,
	}

	for i := int64(0); i < sz; i++ {
		qpsCounter.counters[i] = &Counter{}
	}

	return qpsCounter
}

// AddCount 增加计数
func (q *QPSCounter) AddCount(count int64) {
	if len(q.counters) == 0 {
		panic("QPSCounter not init")
	}

	now := time.Now().Unix()
	index := now % q.size

	q.lock.Lock()
	defer q.lock.Unlock()

	counter := q.counters[index]
	timestamp, total := counter.Serialize()
	if timestamp == now {
		counter.Deserialize(timestamp, total+count)
	} else {
		counter.Deserialize(now, count)
	}
}

// SumCount 统计qps，不加锁，提高效率
func (q *QPSCounter) SumCount() int64 {
	var sum int64
	// 统计数据开始时间
	fromTime := time.Now().Unix() - q.size
	for _, counter := range q.counters {
		timestamp, total := counter.Serialize()
		if timestamp > fromTime {
			sum += total
		}
	}

	return sum / q.size
}
