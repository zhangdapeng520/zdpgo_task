package zdpgo_task

import (
	"time"
)

/*
@Time : 2022/5/7 12:55
@Author : 张大鹏
@File : task
@Software: Goland2021.3.1
@Description: 任务管理核心代码
*/

type Task struct {
	Config *Config // 配置
}

func New() *Task {
	return NewWithConfig(Config{})
}

func NewWithConfig(cfg Config) *Task {
	t := Task{}
	t.Config = &cfg
	return &t
}

// RunTimer 执行定时任务
// @param stopCh 退出通道，用于通知什么时候退出此任务
// @param timerSeconds 定时间隔，每隔多久执行一次任务，单位秒
// @param taskFunc 要执行的任务函数
// @param exitFunc 退出任务之前要执行的函数列表
func (t *Task) RunTimer(stopCh <-chan struct{}, timerSeconds int, taskFunc func(...interface{}), exitFunc ...func()) {
	go func() {
		// 任务执行完毕之后的退出函数
		defer func() {
			if exitFunc != nil && len(exitFunc) > 0 {
				for _, ef := range exitFunc {
					ef()
				}
			}
		}()

		// 定时器
		timer := time.NewTicker(time.Second * time.Duration(timerSeconds))

		// 定时执行任务
		for range timer.C {
			select {
			case <-stopCh:
				return
			default:
				taskFunc()
			}
		}
	}()
	return
}
