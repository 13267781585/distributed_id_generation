package model

type DataSource int64

func (s DataSource) String() string {
	switch s {
	case MysqlSource:
		return "mysql_source"
	case MysqlCacheSource:
		return "mysql_cache_source"
	case RedisSource:
		return "redis_source"
	case RedisCacheSource:
		return "redis_cache_source"
	default:
		return "unknown_source"
	}
}

const (
	MysqlSource DataSource = iota
	RedisSource
	MysqlCacheSource
	RedisCacheSource
)

type Bound struct {
	Start int64
	End   int64
}

type UUIDBound struct {
	IntervalBound *Bound
	Source        DataSource
}
