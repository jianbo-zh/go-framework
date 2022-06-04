package config

import (
	"github.com/jianbo-zh/go-config"
)

func logConfig() {

	// 优先使用命令行参数
	logLevel := config.GetString("log.level", config.GetString("LOG_LEVEL", "info"))
	logFile := config.GetString("log.file", config.GetString("LOG_FILE", "app.log"))

	config.Add("log", config.StrMap{
		"level":   logLevel,                               // 默认日志等级
		"file":    logFile,                                // 日志文件路径
		"format":  config.GetString("LOG_FORMAT", "json"), // 日志格式
		"stdout":  config.GetBool("LOG_STDOUT", true),     // 是否标准输出
		"stderr":  config.GetBool("LOG_STDERR", false),    // 是否错误输出
		"labels":  config.GetString("LOG_LABELS", ""),     // 日志元数据
		"sublvls": config.GetString("LOG_SUBLVLS", ""),    // 子模块日志等级
	})

}
