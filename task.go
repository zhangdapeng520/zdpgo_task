package zdpgo_task

import (
	"context"
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

// RunTimer 执行定时任务，任务会按照指定间隔时间重复执行
// @param stopCh 退出通道，用于通知什么时候退出此任务
// @param timerSeconds 定时间隔，每隔多久执行一次任务，单位毫秒
// @param taskFunc 要执行的任务函数
// @param exitFunc 退出任务之前要执行的函数列表
func (t *Task) RunTimer(stopCh <-chan struct{}, timerMilliSeconds int, taskFunc func(...interface{}),
	exitFunc ...func()) {
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
		timer := time.NewTicker(time.Millisecond * time.Duration(timerMilliSeconds))

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
}

// RunTimeout 运行超时任务，任务会按照指定间隔时间重复执行，超过指定时间后自动退出，也可以主动取消
// @param intervalMilliSecond 间隔时间，单位毫秒
// @param timeoutSecond 超时时间，单位秒
// @param taskFunc 任务函数
// @param taskFunc 任务退出执行要执行的函数列表
// @return context.CancelFunc 取消函数，可以用该函数对象主动取消任务
func (t *Task) RunTimeout(intervalMilliSecond, timeoutSecond int, taskFunc func(...interface{}),
	exitFunc ...func()) context.CancelFunc {
	// 创建的上下文对象和取消函数
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSecond))

	// 执行任务
	go func(ctx2 context.Context, interval int) {
		// 任务执行完毕之后的退出函数
		defer func() {
			if exitFunc != nil && len(exitFunc) > 0 {
				for _, ef := range exitFunc {
					ef()
				}
			}
		}()

		// 定时器
		timer := time.NewTicker(time.Millisecond * time.Duration(interval))

		// 定时执行任务
		for range timer.C {
			select {
			case <-ctx.Done():
				return
			default:
				taskFunc()
			}
		}
	}(ctx, intervalMilliSecond)

	// 将取消函数返回，用于主动取消任务
	return cancel
}
