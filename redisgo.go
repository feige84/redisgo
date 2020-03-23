package redisgo

import (
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

func (rc *RedisInfo) Close() error {
	return rc.p.Close()
}

// associate with config key.
func (rc *RedisInfo) joinPrefix(originKey string) string {
	if rc.key != "" {
		return fmt.Sprintf("%s:%s", rc.key, originKey)
	}
	return originKey
}
