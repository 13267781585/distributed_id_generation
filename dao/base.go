package dao

import "uuid_server/model"

type DaoFetcher interface {
	GetNextInterval(bizCode int64, interval int64) (start int64, end int64, source model.DataSource, err error)
}
