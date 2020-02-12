package redisgo

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	redisGo, err := NewRedisGo(`{"prefix":"doudashi","conn":"127.0.0.1","dbNum":"0","password":"","maxIdle":"0", "maxActive":"", "idleTimeout":"0"}`)
	if err != nil {
		panic(err)
	}

	ok, err := redisGo.HIncrBy("myhash", "field", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(ok)
	//for _, a := range aa {
	//	fmt.Println(string(a.([]byte)))
	//}
}
