package redisgo

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

func (rc *RedisInfo) ClientList() (string, error) {
	c := rc.p.Get()
	defer c.Close()
	return redis.String(c.Do("CLIENT", "LIST"))
}

func (rc *RedisInfo) DbSize() int64 {
	num, _ := redis.Int64(rc.do("DBSIZE"))
	return num
}

func (rc *RedisInfo) Time() int64 {
	timestamp, _ := redis.Int64s(rc.do("TIME"))
	if len(timestamp) > 0 {
		return timestamp[0]
	}
	return time.Now().Unix()
}
