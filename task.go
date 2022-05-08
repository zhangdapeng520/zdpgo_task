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
	BackgroundTaskMap map[string]BackgroundTaskContainer // 后台任务字典
	TimerTaskMap      map[string]TimerTaskContainer      // 定时任务字典
}

// TaskFunc 普通任务函数
type TaskFunc func(...interface{}) (TaskResult, error)

// TaskContainer 普通任务容器
type TaskContainer struct {
	Func TaskFunc // 任务
}

// BackgroundTaskFunc 后台任务方法
type BackgroundTaskFunc func(...interface{})

// BackgroundTaskContainer 后台任务容器
type BackgroundTaskContainer struct {
	Running bool
	Func    func(...interface{}) // 任务
}

// TimerTaskFunc 定时器任务方法
type TimerTaskFunc func(...interface{})

// TimerTaskContainer 定时任务容器
type TimerTaskContainer struct {
	Running      bool                                 // 是否正在运行
	ExitChan     chan bool                            // 控制是否退出的通道
	TimerSeconds int                                  // 监听时间
	Func         func(chan bool, int, ...interface{}) // 任务
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
