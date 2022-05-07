package zdpgo_task

/*
@Time : 2022/5/7 12:56
@Author : 张大鹏
@File : config
@Software: Goland2021.3.1
@Description: 配置相关
*/

type Config struct {
	Debug bool `yaml:"debug" json:"debug" env:"debug"`
}
