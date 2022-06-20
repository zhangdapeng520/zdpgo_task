package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_task/flow"
	"time"
)

/*
@Time : 2022/6/20 9:59
@Author : 张大鹏
@File : main.go
@Software: Goland2021.3.1
@Description:
*/

func main() {
	f1 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("任务1开始了")
		time.Sleep(time.Millisecond * 1000)
		return 1, nil
	}

	f2 := func(r map[string]interface{}) (interface{}, error) {
		time.Sleep(time.Millisecond * 1000)
		fmt.Println("任务2开始了", r["f1"])
		return "some results", nil // errors.New("Some error")
	}

	f3 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("任务3开始了", r["f1"])
		return nil, nil
	}

	f4 := func(r map[string]interface{}) (interface{}, error) {
		fmt.Println("任务4开始了", r)
		return nil, nil
	}

	// 参数1：goroutine名字
	// 参数2：方法参数
	// 参数3：方法
	res, err := flow.New().
		Add("f1", nil, f1).
		Add("f2", []string{"f1"}, f2).
		Add("f3", []string{"f1"}, f3).
		Add("f4", []string{"f2", "f3"}, f4).
		Do()

	fmt.Println(res, err)
}
