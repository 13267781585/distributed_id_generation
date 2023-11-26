package logic

import (
	"fmt"
	"uuid_server/model"
)

const (
	BizNum int64 = 16
)

var MysqlUuidPool *UUIDPool
var RedisUuidPool *UUIDPool

func InitPool() {
	MysqlUuidPool = NewUUIDPool(NewDefaultMysqlConfig())
	RedisUuidPool = NewUUIDPool(NewDefaultRedisConfig())
}

func StopPool() {
	MysqlUuidPool.Stop()
	RedisUuidPool.Stop()
}

func NewUUIDPool(bufferConfig *Config) *UUIDPool {
	pool := &UUIDPool{
		buffers: make(map[int64]*UUIDBuffer),
	}
	for i := int64(0); i < BizNum; i++ {
		bufferConfig.BizCode = i
		buffer := NewUUIDBuffer(bufferConfig)
		buffer.Start()
		pool.buffers[i] = buffer
	}
	return pool
}

type UUIDPool struct {
	buffers map[int64]*UUIDBuffer
}

func (p *UUIDPool) GetUUIDBounds(biz int64, count int64) ([]*model.UUIDBound, error) {
	buffer, ok := p.buffers[biz]
	if !ok {
		return nil, fmt.Errorf("invalid bizCode:%v", biz)
	}

	return buffer.GetUUIDBound(count)
}

func (p *UUIDPool) Stop() {
	for _, buffer := range p.buffers {
		buffer.Stop()
	}
}
