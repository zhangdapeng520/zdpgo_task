package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_task"
	"time"
)

/*
@Time : 2022/6/20 10:16
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/

func main() {
	task := zdpgo_task.New()

	// 释放协程池
	defer task.Close()

	// 定时任务
	stopChan := make(chan struct{})
	task.RunTimer(stopChan, 3000, func() {
		fmt.Println("执行定时任务")
	})

	time.Sleep(time.Second * 10)
	stopChan <- struct{}{} // 停止任务
}
