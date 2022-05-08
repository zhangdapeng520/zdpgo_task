package zdpgo_task

/*
@Time : 2022/5/7 13:37
@Author : 张大鹏
@File : start
@Software: Goland2021.3.1
@Description: 启动任务相关
*/

func (task *Task) Start(taskName string, args ...interface{}) (result TaskResult, err error) {
	if task.TaskMap == nil {
		return
	}
	if taskContainer, ok := task.TaskMap[taskName]; ok {
		result, err = taskContainer.Func(args...)
	}
	return
}

func (task *Task) StartBackground(taskName string, args ...interface{}) {
	if task.BackgroundTaskMap == nil {
		return
	}
	if taskContainer, ok := task.BackgroundTaskMap[taskName]; ok {
		taskContainer.Func(args...)
	}
}

// StartTimer 启动定时任务
func (task *Task) StartTimer(taskName string, args ...interface{}) {
	if task.TimerTaskMap == nil {
		return
	}
	if taskContainer, ok := task.TimerTaskMap[taskName]; ok {
		if !taskContainer.Running {
			if taskContainer.TimerSeconds <= 0 {
				taskContainer.TimerSeconds = 1 // 默认1秒
			}
			go taskContainer.Func(taskContainer.ExitChan, taskContainer.TimerSeconds, args...)
			taskContainer.Running = true
			return
		}
	}
	return
}
