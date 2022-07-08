package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_task"
	"math/rand"
	"time"
)

/*
@Time : 2022/6/20 10:16
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/
var (
	runTimes = 1000
)

func myFunc(arg interface{}) {
	seconds := arg.(int)
	fmt.Println(fmt.Sprintf("模拟一次请求，需要%d秒钟", seconds))
	time.Sleep(time.Second * time.Duration(seconds))
}

func main() {
	task := zdpgo_task.NewWithConfig(&zdpgo_task.Config{
		PoolSize:        100, // 最多同时执行100个任务
		TaskFuncWithArg: myFunc,
	})

	// 释放协程池
	defer task.Close()

	// 提交任务
	fmt.Println("提交任务。。。")
	for i := 0; i < runTimes; i++ {
		seconds := rand.Intn(10)
		task.AddTaskWithArg(seconds)
	}
	task.Wg.Wait()

	// 查看结果
	fmt.Printf("运行中的Goroutine数量： %d\n", task.GetGoroutineNum())
	fmt.Println("任务执行完毕")
}
