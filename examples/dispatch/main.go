package main

import (
	"fmt"
	"github.com/mborders/artifex"
	"time"
)

/*
@Time : 2022/6/17 17:57
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/

func main() {
	// 开启10个goroutine，最多100个任务队列
	d := artifex.NewDispatcher(10, 100)
	d.Start()
	defer d.Stop()

	// 分配100个任务
	for i := 0; i < 1000; i++ {
		err := d.Dispatch(func() {
			// do something
			fmt.Println("do something111")
		})
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Second * 3)
}
