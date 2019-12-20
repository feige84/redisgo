package redisgo

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

func (rc *RedisInfo) Del(key string) error {
	_, err := rc.do("DEL", redis.Args{}.Add(key)...)
	return err
}

func (rc *RedisInfo) DelKeys(pattern string) error {
	c := rc.p.Get()
	defer c.Close()
	keyList, err := rc.Keys(pattern)
	if err != nil {
		return err
	}
	for _, key := range keyList {
		err = c.Send("DEL", key)
		if err != nil {
			return err
		}
	}
	err = c.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (rc *RedisInfo) Expire(key string, life int64) error {
	if _, err := rc.do("EXPIRE", redis.Args{}.Add(key).AddFlat(life)...); err != nil {
		return err
	}
	return nil
}

func (rc *RedisInfo) Exists(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", redis.Args{}.Add(key)...))
	if err != nil {
		return false
	}
	return v
}

func (rc *RedisInfo) Ttl(key string) int64 {
	ttl, _ := redis.Int64(rc.do("TTL", redis.Args{}.Add(key)...))
	return ttl
}

func (rc *RedisInfo) Type(key string) string {
	resultData, _ := redis.String(rc.do("TYPE", redis.Args{}.Add(key)...))
	return resultData
}

func (rc *RedisInfo) Scan(start, pattern string) ([]interface{}, error) {
	if rc == nil {
		return nil, errors.New("redis pool is nil")
	}
	c := rc.p.Get()
	defer c.Close()
	pattern = fmt.Sprintf("%s:%s", rc.key, pattern)
	return redis.Values(c.Do("SCAN", start, "MATCH", pattern, "COUNT", 2000))
}

//OK
func (rc *RedisInfo) Keys(pattern string) ([]string, error) {
	if rc == nil {
		return nil, errors.New("redis pool is nil")
	}
	c := rc.p.Get()
	defer c.Close()

	start := 0
	var err error
	var reply []interface{}
	result := []string{}
	for {
		if pattern != "" {
			reply, err = redis.Values(c.Do("SCAN", start, "MATCH", pattern+"*"))
		} else {
			reply, err = redis.Values(c.Do("SCAN", start))
		}
		if err != nil {
			panic(err)
		}
		if len(reply) > 0 {
			start, _ = strconv.Atoi(string(reply[0].([]byte)))
			if start > 0 {
				list := reply[1].([]interface{})
				if len(list) > 0 {
					for _, v := range list {
						result = append(result, string(v.([]byte)))
					}
				}
			} else {
				break
			}
		}
	}
	return result, err
}
