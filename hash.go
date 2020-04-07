package redisgo

import "github.com/gomodule/redigo/redis"

func (rc *RedisInfo) HExists(key, field interface{}) bool {
	v, err := redis.Bool(rc.do("HEXISTS", redis.Args{}.Add(key).AddFlat(field)...))
	if err != nil {
		return false
	}
	return v
}
func (rc *RedisInfo) HGet(key string, field interface{}) (interface{}, error) {
	return rc.do("HGET", redis.Args{}.Add(key).AddFlat(field)...)
}

func (rc *RedisInfo) HSet(key string, field, value interface{}) (int64, error) {
	return redis.Int64(rc.do("HSET", redis.Args{}.Add(key).AddFlat(field).AddFlat(value)...))
}

func (rc *RedisInfo) HDel(key string, field interface{}) (int64, error) {
	return redis.Int64(rc.do("HDEL", redis.Args{}.Add(key).AddFlat(field)...))
}

func (rc *RedisInfo) HIncrBy(key string, field interface{}, increment int64) (int64, error) {
	return redis.Int64(rc.do("HINCRBY", redis.Args{}.Add(key).AddFlat(field).AddFlat(increment)...))
}

func (rc *RedisInfo) HIncrByFloat(key string, field interface{}, increment int64) (int64, error) {
	return redis.Int64(rc.do("HINCRBYFLOAT", redis.Args{}.Add(key).AddFlat(field).AddFlat(increment)...))
}

func (rc *RedisInfo) HLen(key string) (int64, error) {
	return redis.Int64(rc.do("HLEN", redis.Args{}.Add(key)...))
}

func (rc *RedisInfo) HMGet(key string, subKey1, subKey2 interface{}) ([]interface{}, error) {
	return redis.Values(rc.do("HMGET", redis.Args{}.Add(key).AddFlat(subKey1).AddFlat(subKey2)...))
}

func (rc *RedisInfo) HMSet(key string, s ...interface{}) error {
	if _, err := rc.do("HMSET", redis.Args{}.Add(key).AddFlat(s)...); err != nil {
		return err
	}
	return nil
}

func (rc *RedisInfo) HMGetAll(key string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	keys, err := redis.Values(rc.do("HKEYS", redis.Args{}.Add(key)...))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	values, err := redis.Values(rc.do("HVALS", redis.Args{}.Add(key)...))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	for i, k := range keys {
		if values[i] != nil {
			if val, exists := values[i].([]byte); exists {
				result[string(k.([]byte))] = val
			}
		}
	}
	return result, nil
}

func (rc *RedisInfo) HGetAll(key string) (map[string]string, error) {
	reply, err := redis.StringMap(rc.do("HGETALL", redis.Args{}.Add(key)...))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return reply, nil
}
