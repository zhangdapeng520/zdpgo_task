package zdpgo_task

/*
@Time : 2022/5/7 13:37
@Author : 张大鹏
@File : stop
@Software: Goland2021.3.1
@Description: 停止任务相关
*/

// StopBackground 停止定时器
func (task *Task) StopBackground(taskName string, quit chan bool) {
	if taskName == "" {
		return
	}
	if task.BackgroundTaskMap == nil {
		return
	}
	if _, ok := task.BackgroundTaskMap[taskName]; ok {
		quit <- true // 传入退出信号
	}
}

// StopTimer 停止定时器
func (task *Task) StopTimer(taskName string) {
	if taskName == "" {
		return
	}
	if task.TimerTaskMap == nil {
		return
	}
	if taskContainer, ok := task.TimerTaskMap[taskName]; ok {
		taskContainer.ExitChan <- true // 传入退出信号
		taskContainer.Running = false  // 修改运行状态
	}
}
