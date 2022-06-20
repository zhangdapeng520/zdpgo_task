package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_task/ants"
	"math/rand"
	"sync"
	"time"
)

/*
@Time : 2022/6/20 10:16
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/

func myFunc() {
	seconds := rand.Intn(30)
	fmt.Println(fmt.Sprintf("模拟一次请求，需要%d秒钟", seconds))
	time.Sleep(time.Second * time.Duration(seconds))
}

func main() {
	var (
		wg       sync.WaitGroup
		runTimes = 1000
	)

	// 初始化协程池
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc()
		wg.Done()
	})

	// 释放协程池
	defer p.Release()

	// 提交任务
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Invoke(int32(i))
	}
	wg.Wait()

	// 查看结果
	fmt.Printf("运行中的Goroutine数量： %d\n", p.Running())
	fmt.Println("任务执行完毕")
}
