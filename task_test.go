package zdpgo_task

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

/*
@Time : 2022/5/7 13:40
@Author : 张大鹏
@File : task_test
@Software: Goland2021.3.1
@Description: 任务测试
*/

func getTask() *Task {
	return New()
}

// 测试运行定时任务
func TestTask_RunTimer(t *testing.T) {
	var (
		quit chan struct{}
	)
	task := getTask()

	// 退出通道
	fmt.Println(runtime.NumGoroutine())
	quit = make(chan struct{})

	// 执行定时任务
	fmt.Println(runtime.NumGoroutine())
	task.RunTimer(quit, 1, func(args ...interface{}) {
		fmt.Println("要执行的任务。。。。")
	})
	fmt.Println(runtime.NumGoroutine())

	// 退出定时任务
	go func(quit chan struct{}) {
		time.Sleep(time.Second * 6)
		quit <- struct{}{}
	}(quit)

	// 退出主程序
	time.Sleep(time.Second * 10)
	fmt.Println("main exit")
	fmt.Println(runtime.NumGoroutine())
}
