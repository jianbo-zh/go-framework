package bootstrap

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall"

	"github.com/fvbock/endless"
	"github.com/jianbo-zh/go-config"
	"github.com/jianbo-zh/go-log"
)

// Run 系统引导
func Run() {

	cmdConfig := parseFlag()
	systemInit(cmdConfig)

	startDaemonTask()
	startHttpServer()
}

// systemInit 系统初始化
func systemInit(cmdConfig CmdConfig) {

	// 1. 配置初始化
	initConfig(cmdConfig)

	// 2. 日志配置
	initLogger(LogConf{
		Level:  config.GetString("log.level", "debug"),
		Format: config.GetString("log.format", "color"),
		File:   config.GetString("log.file", "app.log"),
		Stdout: config.GetBool("log.stdout", true),
		Stderr: config.GetBool("log.stderr", true),
	})

	// 3. 初始化翻译
	initTranslation()
}

// startHttpServer 系统启动
func startHttpServer() {

	initRouter()

	server := endless.NewServer(config.GetString("server.host")+":"+config.GetString("server.port"), appRouter)

	server.BeforeBegin = func(add string) {
		ioutil.WriteFile(config.GetString("server.pidfile"), []byte(fmt.Sprintf("%d", syscall.Getpid())), 0666)
	}

	go func() {
		// pprof
		http.ListenAndServe(config.GetString("server.phost")+":"+config.GetString("server.pport"), nil)
	}()

	err := server.ListenAndServe()
	log.Logger("bootstrap").Error(err)
}
