package bootstrap

import (
	allConfig "goframework/cmd/http/config"

	"github.com/jianbo-zh/go-config"
)

func initConfig(conf CmdConfig) {
	// 1. 初始化配置文件及环境变量配置
	config.Setup(config.ConfigOption{
		Type:      conf.CfgType,
		File:      conf.CfgFile,
		EnvPrefix: conf.EnvPrefix,
	})

	// 2. 配置命令行参数
	configCommandLine(conf)

	// 3. 额外指定配置
	allConfig.Setup()
}

func configCommandLine(conf CmdConfig) {

	// 命令行:启动服务配置
	config.Add("server", config.StrMap{
		"host":    conf.Host,
		"port":    conf.Port,
		"pidfile": conf.PidFile,
		"phost":   conf.Phost,
		"pport":   conf.Pport,
	})

	// 命令行:日志配置
	config.Add("log", config.StrMap{
		"level": conf.LogLevel,
		"file":  conf.LogFile,
	})

	// 命令行:存储目录配置
	config.Add("storage", config.StrMap{
		"dir": conf.StorageDir,
	})
}
