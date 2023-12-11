package hook

import (
	"uuid_server/dao"
	"uuid_server/model"
)

type DefaultRedisBoundHook struct {
	uidRedis *dao.UUIDRedis
}

func NewDefaultRedisBoundHook() *DefaultRedisBoundHook {
	return &DefaultRedisBoundHook{
		uidRedis: &dao.UUIDRedis{},
	}
}

func (r *DefaultRedisBoundHook) Store(bounds []*model.Bound, bizCode int64, dataSource model.DataSource) error {
	if len(bounds) == 0 {
		return nil
	}
	return r.uidRedis.StoreBounds(bounds, bizCode, dataSource)
}

func (r *DefaultRedisBoundHook) Load(bizCode int64, dataSource model.DataSource) ([]*model.Bound, error) {
	return r.uidRedis.LoadBounds(bizCode, dataSource)
}
