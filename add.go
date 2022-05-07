package zdpgo_task

import (
	"time"
)

/*
@Time : 2022/5/7 12:58
@Author : 张大鹏
@File : add
@Software: Goland2021.3.1
@Description: 添加任务相关代码
*/

func (task *Task) Add(taskName string, taskFunc TaskFunc) {
	if taskName == "" {
		return
	}
	if task.TaskMap == nil {
		task.TaskMap = make(map[string]TaskContainer)
	}
	task.TaskMap[taskName] = TaskContainer{
		Func: taskFunc,
	}
}

func (task *Task) AddBackground(taskName string, taskFunc BackgroundTaskFunc) {
	if taskName == "" {
		return
	}
	if task.BackgroundTaskMap == nil {
		task.BackgroundTaskMap = make(map[string]BackgroundTaskContainer)
	}
	task.BackgroundTaskMap[taskName] = BackgroundTaskContainer{
		Func: taskFunc,
	}
}

func (task *Task) AddTimer(taskName string, taskFunc TimerTaskFunc) {
	if taskName == "" {
		return
	}
	if task.TimerTaskMap == nil {
		task.TimerTaskMap = make(map[string]TimerTaskContainer)
	}
	task.TimerTaskMap[taskName] = TimerTaskContainer{
		Running:      false,
		ExitChan:     make(chan bool, 1),
		TimerSeconds: 1,
		Func: func(exitChan chan bool, timerSeconds int, args ...interface{}) {
			ticker := time.NewTicker(time.Duration(timerSeconds) * time.Second)
			for range ticker.C {
				select {
				case <-exitChan: // 退出
					ticker.Stop() // 停止定时器
					return
				default:
					taskFunc(args...) // 执行后台任务
				}
			}
		},
	}
}
