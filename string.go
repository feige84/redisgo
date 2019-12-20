package redisgo

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

func (rc *RedisInfo) Get(key string) string {
	resultData, _ := redis.String(rc.do("GET", redis.Args{}.Add(key)...))
	return resultData
}

//存储前先做好数据转存吧，比如json或者xml
func (rc *RedisInfo) Set(key, data interface{}, life int64) error {
	var err error
	var aa interface{}
	if life > 0 {
		aa, err = rc.do("SETEX", redis.Args{}.Add(key).AddFlat(life).AddFlat(data)...)
	} else {
		aa, err = rc.do("SET", redis.Args{}.Add(key).AddFlat(data)...)
	}
	fmt.Println(aa)
	return err
}

//go的gob专用，体积小。
func (rc *RedisInfo) GetBytes(key string, rs interface{}) bool {
	resultData, err := rc.do("GET", redis.Args{}.Add(key)...)
	if err != nil {
		return false
	}
	if resultData == nil {
		return false
	}
	var readBuf bytes.Buffer
	dec := gob.NewDecoder(&readBuf)
	if data, ok := resultData.([]byte); ok {
		_, e := readBuf.Write(data)
		if e != nil {
			return false
		}
	}
	err = dec.Decode(rs)
	if err != nil {
		return false
	}
	return true
	//return resultData
}

//go的gob专用，体积小。
func (rc *RedisInfo) SetBytes(name string, data interface{}, life int64) error {
	var err error
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(data)
	if err != nil {
		panic(err)
	}
	if life > 0 {
		_, err = rc.do("SETEX", name, life, buf.Bytes())
	} else {
		_, err = rc.do("SET", name, buf.Bytes())
	}
	return err
}
func (rc *RedisInfo) Incr(key string) (int64, error) {
	return redis.Int64(rc.do("INCR", redis.Args{}.Add(key)...))
}

// Decr decrease counter in redis.
func (rc *RedisInfo) Decr(key string) (int64, error) {
	return redis.Int64(rc.do("DECR", redis.Args{}.Add(key)...))
}

func (rc *RedisInfo) Incrby(key string, val int) (int64, error) {
	return redis.Int64(rc.do("INCRBY", redis.Args{}.Add(key).AddFlat(val)...))
}

// Decr decrease counter in redis.
func (rc *RedisInfo) Decrby(key string, val int) (int64, error) {
	return redis.Int64(rc.do("DECRBY", redis.Args{}.Add(key).AddFlat(val)...))
}

func (rc *RedisInfo) MGet(keys []string) []interface{} {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.joinPrefix(key))
	}
	values, err := redis.Values(c.Do("MGET", args...))
	if err != nil {
		return nil
	}
	return values
}
