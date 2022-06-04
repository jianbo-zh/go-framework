package bootstrap

import (
	"fmt"

	"github.com/jianbo-zh/go-log"
)

type LogConf struct {
	Level     string            // 默认日志等级
	Format    string            // 日志格式
	File      string            // 日志文件路径
	Stdout    bool              // 是否标准输出
	Stderr    bool              // 是否错误输出
	Labels    map[string]string // 日志元数据
	SubLevels map[string]string // 子模块日志等级
}

func initLogger(conf LogConf) {

	defaultLevel, err := log.LevelFromString(conf.Level)
	if err != nil {
		panic(fmt.Errorf("log level from string error: %w", err))
	}

	labels := conf.Labels
	if labels == nil {
		labels = make(map[string]string)
	}

	var format log.LogFormat
	switch conf.Format {
	case "color":
		format = log.FormatColorizedOutput
	case "nocolor":
		format = log.FormatPlaintextOutput
	case "json":
		format = log.FormatJSONOutput
	default:
		panic(fmt.Sprintf("unrecognized log format '%s'", conf.Format))
	}

	subLvls := make(map[string]log.LogLevel, len(conf.SubLevels))
	for k, v := range conf.SubLevels {
		subLvls[k], err = log.LevelFromString(v)
		if err != nil {
			panic(fmt.Errorf("log level from string error: %w", err))
		}
	}

	log.SetupLogging(log.Config{
		Labels:          labels,
		Format:          format,
		Level:           defaultLevel,
		File:            conf.File,
		Stdout:          conf.Stdout,
		Stderr:          conf.Stderr,
		SubsystemLevels: subLvls,
	})
}
