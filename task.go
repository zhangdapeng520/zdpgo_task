package zdpgo_task

import (
	"context"
	"github.com/zhangdapeng520/zdpgo_log"
	"github.com/zhangdapeng520/zdpgo_task/ants"
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
	Pool   *ants.Pool
}

func New(log *zdpgo_log.Log) *Task {
	return NewWithConfig(&Config{}, log)
}

func NewWithConfig(config *Config, log *zdpgo_log.Log) *Task {
	t := &Task{}

	// 实例化
	t.Log = log
	if config.PoolSize == 0 {
		config.PoolSize = 333
	}
	t.Config = config

	// 任务阻塞：等池中的任务执行完毕了，有空余goroutine可以用了再执行新的任务
	t.Pool, _ = ants.NewPool(config.PoolSize, ants.WithNonblocking(false))

	// 返回
	return t
}

// AddTask 添加任务
func (t *Task) AddTask(taskFunc func()) {
	err := t.Pool.Submit(taskFunc)
	if err != nil {
		t.Log.Error("添加任务失败", "error", err)
	}
}

func (t *Task) Close() {
	t.Pool.Release()
}

func (t Task) GetGoroutineNum() int {
	return t.Pool.Running()
}

// RunTimer 执行定时任务，任务会按照指定间隔时间重复执行
// @param stopCh 退出通道，用于通知什么时候退出此任务
// @param timerSeconds 定时间隔，每隔多久执行一次任务，单位毫秒
// @param taskFunc 要执行的任务函数
// @param exitFunc 退出任务之前要执行的函数列表
func (t *Task) RunTimer(stopCh <-chan struct{}, timerMilliSeconds int, taskFunc func(...interface{}), exitFunc ...func()) {
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
func (t *Task) RunTimerTimeout(intervalMilliSecond, timeoutSecond int, taskFunc func(...interface{}), exitFunc ...func()) context.CancelFunc {
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
	// 创建通道
	ch := make(chan bool, 1)
	defer close(ch)

	// 执行任务
	go func() {
		// 处理通道关闭触发的异常
		defer func() {
			if r := recover(); r != nil {
				t.Log.Debug("向已关闭的通道写入数据", "error", r)
			}
		}()

		// 执行任务
		taskFunc()
		ch <- true
	}()

	// 创建定时器，当函数返回时，它所使用的所有通道都已被清除。
	timer := time.NewTimer(time.Duration(timeoutSecond) * time.Second)
	defer timer.Stop()

	// 监听任务结束情况
	select {
	case <-ch:
		taskFinishFunc()
	case <-timer.C: // 监听到定时器（已超时）
		timeoutFunc()
	}
}
