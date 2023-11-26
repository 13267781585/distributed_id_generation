package dao

import (
	"sync/atomic"
	"uuid_server/model"
)

var redisCount = int64(0)

type UIDCounterRedis struct {
}

func (m *UIDCounterRedis) GetNextInterval(bizCode int64, interval int64) (int64, int64, model.DataSource, error) {
	end := atomic.AddInt64(&redisCount, interval)
	return end - interval, end, model.RedisSource, nil
}
