package zdpgo_task

import (
	"fmt"
	"testing"
	"time"
)

/*
@Time : 2022/5/7 13:40
@Author : 张大鹏
@File : add_test
@Software: Goland2021.3.1
@Description:
*/

func TestTask_Add(t *testing.T) {
	task := getTask()
	task.Add("test1", func(i ...interface{}) (TaskResult, error) {
		fmt.Println(i...)
		return TaskResult{}, nil
	})
	task.Start("test1", 1, 2, 3, 4)

	// 加法
	task.Add("test2", func(args ...interface{}) (TaskResult, error) {
		a := args[0].(int)
		b := args[1].(int)
		return TaskResult{Value: a + b}, nil
	})
	result, err := task.Start("test2", 1, 2)
	if err != nil {
		panic(err)
	}
	value := result.Value.(int)
	fmt.Println(value)
}

func TestTask_AddBackground(t *testing.T) {
	task := getTask()
	f := func(ch chan interface{}, i ...interface{}) {
		for j := 0; j < 10; j++ {
			fmt.Println(i...)
			time.Sleep(time.Second)
		}
		ch <- true
	}
	task.AddBackground("test1", f)
	ch := make(chan interface{}, 1)
	task.StartBackground("test1", ch, 1, 2, 3, 4)
	<-ch
}

func TestTask_AddTimer(t *testing.T) {
	task := getTask()

	// 后台任务
	f := func(i ...interface{}) {
		for {
			fmt.Println("接收到的参数：", i)
			time.Sleep(time.Second)
		}
	}

	// 添加后台任务
	task.AddTimer("test1", f)

	// 执行后台任务
	task.StartTimer("test1", 1, 2, 3, 4)

	// 3秒中以后停止后台任务
	time.Sleep(time.Second * 3)
	fmt.Println("准备停止后台任务。。。")
	task.StopTimer("test1")
}
