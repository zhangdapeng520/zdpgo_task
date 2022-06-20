package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_task/ants"
	"sync"
	"sync/atomic"
)

/*
@Time : 2022/6/20 10:16
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/

var sum int32 // 总和

func myFunc(i interface{}) {
	n := i.(int32)           // 将参数转换为真实的类型
	atomic.AddInt32(&sum, n) // 原子执行相加操作
	fmt.Printf("加上： %d\n", n)
}

func main() {
	var (
		wg       sync.WaitGroup
		runTimes = 1000
	)

	// 初始化协程池
	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) {
		myFunc(i)
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
	fmt.Printf("任务执行完毕，结果是： %d\n", sum)
}
