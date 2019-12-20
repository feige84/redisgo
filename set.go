package redisgo

import "github.com/gomodule/redigo/redis"

func (rc *RedisInfo) SAdd(key interface{}, s ...interface{}) (int64, error) {
	return redis.Int64(rc.do("SADD", redis.Args{}.Add(key).AddFlat(s)...))
}

func (rc *RedisInfo) SPop(key interface{}) (interface{}, error) {
	if reply, err := rc.do("SPOP", redis.Args{}.Add(key)...); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}

func (rc *RedisInfo) SRem(key interface{}, s ...interface{}) (int64, error) {
	return redis.Int64(rc.do("SREM", redis.Args{}.Add(key).AddFlat(s)...))
}

func (rc *RedisInfo) SCard(key interface{}) (int64, error) {
	return redis.Int64(rc.do("SCARD", redis.Args{}.Add(key)...))
}

func (rc *RedisInfo) SIsMember(key, s interface{}) bool {
	v, err := redis.Bool(rc.do("SISMEMBER", redis.Args{}.Add(key).AddFlat(s)...))
	if err != nil {
		return false
	}
	return v
}
