package service

import (
	"uuid_server/logic"
	"uuid_server/model"
)

func NewUUIDService() *UUIDService {
	return &UUIDService{}
}

type UUIDService struct {
}

func (s *UUIDService) GetUUIDBounds(biz int64, count int64) ([]*model.UUIDBound, error) {
	bounds, err := logic.MysqlUuidPool.GetUUIDBounds(biz, count)
	if err == nil {
		return bounds, nil
	}

	bounds, err = logic.RedisUuidPool.GetUUIDBounds(biz, count)
	return bounds, err
}
