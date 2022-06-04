package task

import (
	"goframework/cmd/http/task/cleantask"
	"goframework/cmd/http/task/gencidtask"
	"goframework/cmd/http/task/upqueuetask"
)

// 任务引导
func Setup() {
	// 同步任务
	upqueuetask.Setup()

	// 清理过期文件
	cleantask.Setup()

	// 生成cid任务
	gencidtask.Setup()
}
