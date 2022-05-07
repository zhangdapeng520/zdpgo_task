package zdpgo_task

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
		Running: false,
		Func:    taskFunc,
	}
}
