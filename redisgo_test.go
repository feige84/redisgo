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
	for i := 0; i < 20; i++ {

		fmt.Println(redisGo.Time())
	}
}
