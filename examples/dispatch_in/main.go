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
	// 10 workers, 100 max in job queue
	d := artifex.NewDispatcher(10, 100)
	d.Start()
	defer d.Stop()

	for i := 0; i < 1000; i++ {
		err := d.DispatchIn(func() {
			// do something in 500ms
			fmt.Println("do something", i)
		}, time.Millisecond*500)
		if err != nil {
			panic(err)
		}
	}

	time.Sleep(time.Second * 10)
}
