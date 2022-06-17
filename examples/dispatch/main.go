package main

import (
	"fmt"
	"github.com/mborders/artifex"
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
	err := d.Dispatch(func() {
		// do something
		fmt.Println("do something111")
	})
	if err != nil {
		panic(err)
	}
}
