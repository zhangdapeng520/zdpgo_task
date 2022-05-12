package zdpgo_task

import (
	"context"
	"github.com/zhangdapeng520/zdpgo_log"
	"sync"
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
	Config *Config        // 配置
	Log    *zdpgo_log.Log // 日志
}

func New() *Task {
	return NewWithConfig(Config{})
}

func NewWithConfig(config Config) *Task {
	t := Task{}

	// 创建日志
	if config.LogFilePath == "" {
		config.LogFilePath = "logs/zdpgo/zdpgo_task.log"
	}
	logConfig := zdpgo_log.Config{
		Debug:         false,
		LogLevel:      "",
		IsWriteDebug:  false,
		IsShowConsole: false,
		OpenJsonLog:   false,
		OpenFileName:  false,
		LogFilePath:   "",
		MaxSize:       0,
		MaxBackups:    0,
		MaxAge:        0,
		Compress:      false,
	}
	if config.Debug {
		logConfig.IsShowConsole = true
	}
	t.Log = zdpgo_log.NewWithConfig(logConfig)
	t.Config = &config
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

// RunTimerTimeout 运行间隔超时任务，任务会按照指定间隔时间重复执行，超过指定时间后自动退出，也可以主动取消
// @param intervalMilliSecond 间隔时间，单位毫秒
// @param timeoutSecond 超时时间，单位秒
// @param taskFunc 任务函数
// @param taskFunc 任务退出执行要执行的函数列表
// @return context.CancelFunc 取消函数，可以用该函数对象主动取消任务
func (t *Task) RunTimerTimeout(intervalMilliSecond, timeoutSecond int, taskFunc func(...interface{}),
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

// RunWaitTasks 同时执行多个任务
// @param taskFuncList 任务列表
// @param exitFuncList 任务执行完毕后要执行的函数列表
func (t *Task) RunWaitTasks(taskFuncList []func(), exitFuncList ...func()) {
	// 创建等待组
	var wg = new(sync.WaitGroup)

	// 执行任务列表
	for _, taskFunc := range taskFuncList {
		wg.Add(1)
		go func(wg *sync.WaitGroup, taskFunc func()) {
			defer wg.Done()
			taskFunc() // 执行核心函数
		}(wg, taskFunc)
	}

	// 等待所有任务结束
	wg.Wait()

	// 执行退出函数
	if exitFuncList != nil && len(exitFuncList) > 0 {
		for _, ef := range exitFuncList {
			ef()
		}
	}
}

// RunExitTimeout 运行退出超时任务，达到超时时间后自动结束该任务
// @param timeoutSecond 超时时间
// @param taskFunc 要执行的任务方法
// @param taskFinishFunc 任务正常结束执行的方法
// @param timeoutFunc 任务超时结束执行的方法
func (t *Task) RunExitTimeout(timeoutSecond int, taskFunc func(), taskFinishFunc func(), timeoutFunc func()) {
	// 创建一个管道
	ch := make(chan struct{}, 1)

	// 异步执行任务，执行完毕发送通知
	go func(taskFunc func()) {
		taskFunc()
		ch <- struct{}{}
	}(taskFunc)

	// 监听任务
	select {
	case <-ch: // 正常结束
		taskFinishFunc()
	case <-time.After(time.Duration(timeoutSecond) * time.Second): // 超时结束
		timeoutFunc()
	}
}
