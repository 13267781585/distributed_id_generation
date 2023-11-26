package dao

import (
	"fmt"
	"sync/atomic"
	"uuid_server/model"
)

var mysqlCount = int64(0)

type UIDCounterMysql struct {
}

func (m *UIDCounterMysql) GetNextInterval(bizCode int64, interval int64) (int64, int64, model.DataSource, error) {
	end := atomic.AddInt64(&mysqlCount, interval)
	fmt.Printf("mysql add interval:%v now:%v \n", interval, end)
	return end - interval, end, model.MysqlSource, nil
}
