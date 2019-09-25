package redisgo

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis" //Redis adapter use the redigo client.
	"strconv"
	"strings"
	"time"
)

var (
	DefaultPrefix = "redisgo"
)

type RedisInfo struct {
	p           *redis.Pool
	host        string
	dbNum       int
	password    string
	key         string
	maxIdle     int
	maxActive   int
	idleTimeout int
}

//var (
//	RedisRead  = &RedisInfo{} //业务读取
//	RedisWrite = &RedisInfo{} //业务写入
//)

//func init() {
//【写】
// RedisWrite = NewRedisCache(KEY_PREFIX, 0).Start()
//【读】
// RedisRead = RedisWrite
//}

//config is like {"prefix":"collection key","conn":"connection info","dbNum":"0","password":"","maxIdle":"0", "maxActive":"", "idleTimeout":"0"}
func NewRedisGo(config string) (*RedisInfo, error) {
	var cf map[string]string
	err := json.Unmarshal([]byte(config), &cf)
	if err != nil {
		return nil, err
	}
	if _, ok := cf["prefix"]; !ok {
		cf["prefix"] = DefaultPrefix
	}
	if _, ok := cf["conn"]; !ok {
		return nil, errors.New("config has no conn key")
	}
	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "5000"
	}
	if _, ok := cf["maxActive"]; !ok {
		cf["maxActive"] = "5000"
	}
	if _, ok := cf["idleTimeout"]; !ok {
		cf["idleTimeout"] = "0"
	}

	rc := &RedisInfo{}
	rc.key = cf["prefix"]
	rc.host = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]
	rc.maxIdle, _ = strconv.Atoi(cf["maxIdle"])
	rc.maxActive, _ = strconv.Atoi(cf["maxActive"])
	rc.idleTimeout, _ = strconv.Atoi(cf["idleTimeout"])
	if rc.idleTimeout <= 0 {
		rc.idleTimeout = 3
	}

	//fmt.Printf("%+v\n", rc)

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	if c.Err() != nil {
		return rc, c.Err()
	}

	return rc, nil
}

func (rc *RedisInfo) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		if !strings.Contains(rc.host, ":") {
			rc.host = fmt.Sprintf("%s:%d", rc.host, 6379)
		}
		c, err = redis.Dial("tcp", rc.host,
			redis.DialPassword(rc.password),
			redis.DialDatabase(rc.dbNum),
			redis.DialConnectTimeout(time.Duration(rc.idleTimeout)*time.Second),
			redis.DialReadTimeout(time.Duration(rc.idleTimeout)*time.Second),
			redis.DialWriteTimeout(time.Duration(rc.idleTimeout)*time.Second),
		)
		if err != nil {
			return nil, err
		}
		return
	}

	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:   rc.maxIdle,
		MaxActive: rc.maxActive,
		// 定义链接的超时时间，每次p.Get()的时候会检测这个连接是否超时（超时会关闭，并释放可用连接数）.0就是不超时。传说中的长连接？
		// 文章上说。IdleTimeout 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用。，但是如果时间断会频繁的select切换库。但是redis跟客户端好像没有太大联系。
		// 综上可知，想提高redis的处理速度，可以把idle设置的大一点，当然active的值是一定要比idle大的（0表示不限制）。这里还有一点需要注意的redis本身支持的连接数设置问题，刚才分析的都是client端的情况，如果server只支持100的连接数，那客户端的pool设定再多也没有用，redis的连接数配置又是另外一个话题，这里就不展开讲。
		// MaxActive 可以把MaxActive调大（一般设置为500，1000问题都不大。）但如果redis服务器负载已经很高了（可以看redis-server CPU占用），去调大MaxActive就没多大意义。还是需要根据实际情况来权衡。
		// Wait 可以把Wait设置为true。wait的话必然会加大响应，如果对响应时间要求较高的话，还得从别的途径来解决。

		//如果把这个参数设置0。是不是意味着长连接啦？只要有db的NUM就够了？
		IdleTimeout: time.Duration(rc.idleTimeout) * time.Second,
		Wait:        true,
		Dial:        dialFunc,
	}
}

// actually do the redis cmds, args[0] must be the key name.
func (rc *RedisInfo) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if rc == nil {
		return nil, errors.New("redis pool is nil")
	}
	//if len(args) < 1 {
	//	return nil, errors.New("missing required arguments")
	//}
	if len(args) > 0 {
		args[0] = rc.joinPrefix(args[0].(string))
	}
	c := rc.p.Get()
	defer c.Close()

	//c.Do("INCRBY", "redisDo", 1)
	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *RedisInfo) joinPrefix(originKey string) string {
	if rc.key != "" {
		return fmt.Sprintf("%s:%s", rc.key, originKey)
	}
	return originKey
}

func (rc *RedisInfo) ClientList() (string, error) {
	c := rc.p.Get()
	defer c.Close()
	return redis.String(c.Do("CLIENT", "LIST"))
}

func (rc *RedisInfo) DbSize() int64 {
	num, _ := redis.Int64(rc.do("DBSIZE"))
	return num
}

func (rc *RedisInfo) Type(name string) string {
	resultData, _ := redis.String(rc.do("TYPE", name))
	return resultData
}

func (rc *RedisInfo) Get(name string) string {
	resultData, _ := redis.String(rc.do("GET", name))
	return resultData
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

/*
func (rc *RedisInfo) Set(name string, data interface{}, life int64) {
	jsonData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		panic(jsonErr.Error())
	}
	var err error
	if life > 0 {
		_, err = rc.do("SETEX", name, life, string(jsonData))
	} else {
		_, err = rc.do("SET", name, string(jsonData))
	}
	if err != nil {
		panic(err.Error())
	}
}
*/

//存储前先做好数据转存吧，比如json或者xml
func (rc *RedisInfo) Set(name, data string, life int64) error {
	var err error
	if life > 0 {
		_, err = rc.do("SETEX", name, life, data)
	} else {
		_, err = rc.do("SET", name, data)
	}
	return err
}

//go的gob专用，体积小。
func (rc *RedisInfo) GetBytes(name string, rs interface{}) bool {
	resultData, err := rc.do("GET", name)
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
		panic(err.Error())
	}
	if life > 0 {
		_, err = rc.do("SETEX", name, life, buf.Bytes())
	} else {
		_, err = rc.do("SET", name, buf.Bytes())
	}
	return err
}

//发送订阅，返回接收到的数量
func (rc *RedisInfo) Publish(channel, message string) int64 {
	c := rc.p.Get()
	defer c.Close()
	reply, _ := redis.Int64(c.Do("PUBLISH", channel, message))
	return reply
}

func (rc *RedisInfo) Ttl(key string) int64 {
	ttl, _ := redis.Int64(rc.do("TTL", key))
	return ttl
}

func (rc *RedisInfo) Time() int64 {
	timestamp, _ := redis.Int64s(rc.do("TIME"))
	if len(timestamp) > 0 {
		return timestamp[0]
	}
	return time.Now().Unix()
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
			panic(err.Error())
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

//OK
func (rc *RedisInfo) Del(name string) error {
	_, err := rc.do("DEL", name)
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

func (rc *RedisInfo) LPush(key, content string) int64 {
	num, _ := redis.Int64(rc.do("LPUSH", key, content))
	return num
}

func (rc *RedisInfo) LPop(key string) (interface{}, error) {
	if reply, err := rc.do("LPOP", key); err != nil {
		return nil, err
	} else {
		return reply, nil
	}
}

func (rc *RedisInfo) HMGet(key, subKey1, subKey2 string) ([]interface{}, error) {
	return redis.Values(rc.do("HMGET", key, subKey1, subKey2))
}

func (rc *RedisInfo) HMGetAll(key string) map[string]interface{} {
	result := make(map[string]interface{})
	keys, err := redis.Values(rc.do("HKEYS", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil
		} else {
			panic(err.Error())
		}
	}
	values, err := redis.Values(rc.do("HVALS", key))
	if err != nil {
		if err == redis.ErrNil {
			return nil
		} else {
			panic(err.Error())
		}
	}
	for i, v := range keys {
		if _, exists := values[i].([]byte); exists {
			result[string(v.([]byte))] = string(values[i].([]byte))
		}
	}
	return result
}

func (rc *RedisInfo) HMSet(key string, s ...interface{}) error {
	if _, err := rc.do("HMSET", redis.Args{}.Add(key).AddFlat(s)...); err != nil {
		return err
	}
	return nil
}

func (rc *RedisInfo) Expire(key string, life int64) error {
	if _, err := rc.do("EXPIRE", key, life); err != nil {
		return err
	}
	return nil
}

func (rc *RedisInfo) Exists(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (rc *RedisInfo) Incr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (rc *RedisInfo) Decr(key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

func (rc *RedisInfo) SAdd(key string, s ...interface{}) (int64, error) {
	return redis.Int64(rc.do("SADD", redis.Args{}.Add(key).AddFlat(s)...))
}

func (rc *RedisInfo) SRem(key string, s ...interface{}) (int64, error) {
	return redis.Int64(rc.do("SREM", redis.Args{}.Add(key).AddFlat(s)...))
}

//这部分是做分布式订阅用的。不通用。
type SubscribeMsg struct {
	SendName string `json:"send_name"`
	FuncName string `json:"func_name"`
	FuncId   int64  `json:"func_id"`
	FuncArg  string `json:"func_arg"`
}

func (rc *RedisInfo) Notify(notifyKey, hostname string, funcId int64, funcName, funcArg string) int64 {
	msg := SubscribeMsg{}
	msg.SendName = hostname
	msg.FuncId = funcId
	msg.FuncName = funcName
	msg.FuncArg = funcArg
	jsonData, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return rc.Publish(notifyKey, string(jsonData))
}
