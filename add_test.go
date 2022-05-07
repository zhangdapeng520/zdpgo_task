package zdpgo_task

import (
	"fmt"
	"testing"
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
