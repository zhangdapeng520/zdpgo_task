package zdpgo_task

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/zhangdapeng520/zdpgo_log"
)

/*
@Time : 2022/5/7 13:40
@Author : 张大鹏
@File : task_test
@Software: Goland2021.3.1
@Description: 任务测试
*/

func getTask() *Task {
	return New(zdpgo_log.Tmp)
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
	task.RunTimer(quit, 1000, func(args ...interface{}) {
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

// 测试运行超时任务
func TestTask_RunTimerTimeout(t *testing.T) {
	task := getTask()
	var cancel context.CancelFunc
	cancel = task.RunTimerTimeout(500, 7, func(args ...interface{}) {
		fmt.Println("任务执行中。。。")
	})
	fmt.Println(runtime.NumGoroutine())

	time.Sleep(time.Second * 6)
	cancel()

	time.Sleep(time.Second)
	fmt.Println(runtime.NumGoroutine())
}

// 测试运行多个任务
func TestTask_RunWaitTasks(t *testing.T) {
	task := getTask()
	var funcs []func()
	funcs = append(funcs, func() {
		time.Sleep(time.Second * 3)
		fmt.Println("函数1")
	})
	funcs = append(funcs, func() {
		time.Sleep(time.Second * 2)
		fmt.Println("函数2")
	})
	funcs = append(funcs, func() {
		time.Sleep(time.Second * 1)
		fmt.Println("函数3")
	})

	task.RunWaitTasks(funcs)
}

// 测试任务超时结束
func TestTask_RunExitTimeout(t *testing.T) {
	task := getTask()

	// 普通任务
	task.RunExitTimeout(3, func() {
		fmt.Println("正常任务")
	}, func() {
		fmt.Println("任务正常结束")
	}, func() {
		fmt.Println("任务超时技术")
	})
	fmt.Println("====================")

	// 任务超时
	task.RunExitTimeout(3, func() {
		fmt.Println("正常任务")
		time.Sleep(time.Second * 4)
	}, func() {
		fmt.Println("任务正常结束")
	}, func() {
		fmt.Println("任务超时结束")
	})
	fmt.Println("====================")

	// 任务循环超时
	task.RunExitTimeout(3, func() {
		fmt.Println("正常任务")
		for {
			time.Sleep(time.Second * 2)
		}
	}, func() {
		fmt.Println("任务正常结束")
	}, func() {
		fmt.Println("任务超时结束")
	})
	fmt.Println("====================")

	// 校验任务是否正常结束的布尔值
	tmpFlag := make(chan bool, 1)
	task.RunExitTimeout(3, func() {
		fmt.Println("正常任务")
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second * 2)
		}
	}, func() {
		fmt.Println("任务正常结束")
		tmpFlag <- true
	}, func() {
		fmt.Println("任务超时结束")
		tmpFlag <- false
	})
	fmt.Println(<-tmpFlag) // 取出结果
}

// 测试任务超时结束
func TestTask_AddTaskWithArg(t *testing.T) {
	task := NewWithConfig(&Config{
		PoolSize:        100,
		TaskFuncWithArg: func(arg interface{}) {},
	}, zdpgo_log.Tmp)

	// 普通任务
	for i := 0; i < 100000; i++ {
		task.AddTaskWithArg(i)
	}
	task.Wg.Wait()
}
