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
		if !taskContainer.Running {
			result, err = taskContainer.Func(args...)
			return
		}
	}
	return
}
