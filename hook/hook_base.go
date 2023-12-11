package hook

import "uuid_server/model"

type IntervalBoundHook interface {
	Store(bounds []*model.Bound, bizCode int64, dataSource model.DataSource) error
	Load(bizCode int64, dataSource model.DataSource) ([]*model.Bound, error)
}
