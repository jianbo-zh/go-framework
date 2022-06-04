package gencidtask

// Setup 启动任务
func Setup() {
	setupListenTask()
	setupWorkerTask()
	setupSmallWorkerPool()
}

// setupListenTask 监听大文件处理队列
func setupListenTask() {
	// 大文件 push 到大文件队列里面
	go largeFileListen()

	// 普通文件 push 到普通文件队列里面
	go normalFileListen()

	// 小文件 push 到小文件队列里面
	go smallFileListen()
}

// setupWorkerTask 处理大文件任务
func setupWorkerTask() {
	// 启动大文件处理器
	go startLargeFileWorker()

	// 启动普通文件处理器
	go startNormalFileWorker()

	// 启动小文件文件处理器
	go startSmallFileWorker()
}

func setupSmallWorkerPool() {
	// 创建小文件工作池
	go createSmallWorkerPool()
}
