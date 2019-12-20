package redisgo

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
)

//发送订阅，返回接收到的数量
func (rc *RedisInfo) Publish(channel, message string) int64 {
	c := rc.p.Get()
	defer c.Close()
	reply, _ := redis.Int64(c.Do("PUBLISH", channel, message))
	return reply
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
