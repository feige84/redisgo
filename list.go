package redisgo

import "github.com/gomodule/redigo/redis"

func (rc *RedisInfo) LLen(key string) int64 {
	num, _ := redis.Int64(rc.do("LLEN", redis.Args{}.Add(key)...))
	return num
}

func (rc *RedisInfo) LRange(key string, start, end int) ([]interface{}, error) {
	return redis.Values(rc.do("LRANGE", redis.Args{}.Add(key).AddFlat(start).AddFlat(end)...))
}

func (rc *RedisInfo) LPush(key string, content ...interface{}) int64 {
	num, _ := redis.Int64(rc.do("LPUSH", redis.Args{}.Add(key).AddFlat(content)...))
	return num
}

func (rc *RedisInfo) LPop(key string) (interface{}, error) {
	if reply, err := rc.do("LPOP", redis.Args{}.Add(key)...); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}

func (rc *RedisInfo) RPush(key string, content ...interface{}) int64 {
	num, _ := redis.Int64(rc.do("RPUSH", redis.Args{}.Add(key).AddFlat(content)...))
	return num
}

func (rc *RedisInfo) RPop(key string) (interface{}, error) {
	if reply, err := rc.do("RPOP", redis.Args{}.Add(key)...); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}
