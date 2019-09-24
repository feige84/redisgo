package redisgo

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	redisGo, err := NewRedisGo(`{"prefix":"redisgo","conn":"127.0.0.1","dbNum":"11","password":"","maxIdle":"0", "maxActive":"", "idleTimeout":"0"}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(redisGo.Get("user_count_2772"))
}
