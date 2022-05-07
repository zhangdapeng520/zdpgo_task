package zdpgo_task

/*
@Time : 2022/5/7 12:55
@Author : 张大鹏
@File : task
@Software: Goland2021.3.1
@Description: 任务管理核心代码
*/

type Task struct {
	Config            *Config                            // 配置
	TaskMap           map[string]TaskContainer           // 任务字典
	BackgroundTaskMap map[string]BackgroundTaskContainer // 任务字典
}
type TaskFunc func(...interface{}) (TaskResult, error)
type TaskContainer struct {
	Running bool     // 是否正在运行
	Func    TaskFunc // 任务
}
type BackgroundTaskFunc func(chan interface{}, ...interface{})
type BackgroundTaskContainer struct {
	Running bool               // 是否正在运行
	Func    BackgroundTaskFunc // 任务
}

type TaskResult struct {
	Value interface{}
}

func New() *Task {
	t := Task{}
	return &t
}

func NewWithConfig(cfg Config) *Task {
	t := Task{}
	t.Config = &cfg
	return &t
}
