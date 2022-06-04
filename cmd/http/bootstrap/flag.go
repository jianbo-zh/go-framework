package bootstrap

import "flag"

type CmdConfig struct {
	Host       string // 监听地址
	Port       string // 监听端口号
	Phost      string // pprof 监听地址
	Pport      string // pport 监听端口
	PidFile    string // PID文件路径
	CfgType    string // 配置文件类型
	CfgFile    string // 配置文件路径
	EnvPrefix  string // 环境变量前缀
	LogFile    string // 日志文件路径
	LogLevel   string // 日志记录等级
	StorageDir string // 数据存储目录
}

func parseFlag() CmdConfig {
	var initConf CmdConfig

	flag.StringVar(&initConf.PidFile, "pid", "app.pid", "save pid file path")

	flag.StringVar(&initConf.Host, "host", "0.0.0.0", "listen address")
	flag.StringVar(&initConf.Port, "port", "8080", "listen port")

	flag.StringVar(&initConf.Phost, "phost", "127.0.0.1", "pprof listen address")
	flag.StringVar(&initConf.Pport, "pport", "8081", "pprof listen port")

	flag.StringVar(&initConf.CfgType, "conf.type", "env", "config file type")
	flag.StringVar(&initConf.CfgFile, "conf.file", ".env", "config file path")
	flag.StringVar(&initConf.EnvPrefix, "conf.envPrefix", "", "config env prefix")

	flag.StringVar(&initConf.LogFile, "log.file", "", "log file path")
	flag.StringVar(&initConf.LogLevel, "log.level", "", "log level")

	flag.StringVar(&initConf.StorageDir, "storage.dir", "", "data storage dir")

	flag.Parse()

	return initConf
}
