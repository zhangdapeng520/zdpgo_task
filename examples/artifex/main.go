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

	d.Dispatch(func() {
		// do something
		fmt.Println("do something")
	})

	err := d.DispatchIn(func() {
		// do something in 500ms
		fmt.Println("do something")
	}, time.Millisecond*500)
	if err != nil {
		panic(err)
	}

	// Returns a DispatchTicker
	dt, err := d.DispatchEvery(func() {
		// do something every 250ms
		fmt.Println("do something")
	}, time.Millisecond*250)

	// Stop a given DispatchTicker
	dt.Stop()

	// Returns a DispatchCron
	dc, err := d.DispatchCron(func() {
		// do something every 1s
	}, "*/1 * * * * *")

	// Stop a given DispatchCron
	dc.Stop()

	// Stop a dispatcher and all its workers/tickers
	d.Stop()
}
