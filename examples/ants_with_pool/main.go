package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_task/ants"
	"math/rand"
	"sync"
	"time"
)

/*
@Time : 2022/6/20 10:08
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description: 计算大量整数和的程序
*/

func demoFunc() {
	seconds := rand.Intn(30)
	fmt.Println(fmt.Sprintf("模拟一次请求，需要%d秒钟", seconds))
	time.Sleep(time.Second * time.Duration(seconds))
}

func main() {
	// 释放ants的默认协程池
	defer ants.Release()

	var wg sync.WaitGroup
	// 任务函数
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		_ = ants.Submit(syncCalculateSum) // 提交任务到默认协程池
	}
	wg.Wait()

	// 查看输出
	fmt.Printf("运行中的Gotouine数量: %d\n", ants.Running())
	fmt.Printf("任务执行完毕\n")
}
