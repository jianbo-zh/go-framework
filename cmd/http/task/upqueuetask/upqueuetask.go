package upqueuetask

import (
	"goframework/cmd/http/service/filesvc"

	"github.com/jianbo-zh/go-log"
)

var taskLog = log.Logger("task")

func Setup() {
	UpdateQueueTaskPending()
}

func UpdateQueueTaskPending() {
	err := filesvc.UpdateQueueHashPending()
	if err != nil {
		taskLog.Errorf("update queue hash error: %s", err.Error())
		return
	}

	taskLog.Infof("update queue hash pending success")
}
