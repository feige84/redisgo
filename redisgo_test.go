package redisgo

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	redisGo, err := NewRedisGo(`{"prefix":"dds_update","conn":"127.0.0.1","dbNum":"11","password":"","maxIdle":"0", "maxActive":"", "idleTimeout":"0"}`)
	if err != nil {
		panic(err)
	}
	aa, err := redisGo.LRange("proxy_list", 0, -1)
	if err != nil {
		panic(err)
	}
	for _, a := range aa {
		fmt.Println(string(a.([]byte)))
	}
}
