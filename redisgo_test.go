package redisgo

import (
	"fmt"
	"testing"
)

type DyAwemeQueue struct {
	QId       int64  `json:"q_id"`       // ID
	QStatus   int64  `json:"q_status"`   // 1，未执行，2,正在执行，3，出错
	QDateline int64  `json:"q_dateline"` // 生成队列时间
	QRunTime  int64  `json:"q_run_time"` // 队列运行时间
	QPage     int64  `json:"q_page"`     //已跑页数
	QCursor   int64  `json:"q_cursor"`   //下次请求的cursor
	QCode     string `json:"q_code"`     // 错误码
	QMsg      string `json:"q_msg"`      // 错误消息
}

func TestExecute(t *testing.T) {
	redisGo, err := NewRedisGo(`{"prefix":"doudashi","conn":"127.0.0.1","dbNum":"11","password":"","maxIdle":"0", "maxActive":"", "idleTimeout":"0"}`)
	if err != nil {
		panic(err)
	}

	aa, err := redisGo.HDel("user_queue_hash", 2367951)
	//queue := DyAwemeQueue{
	//	QId:       111111,
	//	QStatus:   0,
	//	QDateline: 0,
	//	QRunTime:  0,
	//	QPage:     3,
	//	QCursor:   22222,
	//	QCode:     "",
	//	QMsg:      "",
	//}
	//jsonData,err:=json.Marshal(queue)
	//if err != nil {
	//	panic(err)
	//}
	//aa, err := redisGo.HSet("user_aweme_queue", "22222", string(jsonData))
	//aa, err := redisGo.HMGetAll("user_aweme_queue")
	if err != nil {
		panic(err)
	}
	fmt.Println(aa)
	//for _, a := range aa {
	//	fmt.Println(string(a.([]byte)))
	//}
}
