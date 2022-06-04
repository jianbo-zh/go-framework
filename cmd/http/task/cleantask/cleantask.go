package cleantask

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jianbo-zh/go-config"
	"github.com/jianbo-zh/go-log"
)

var taskLog = log.Logger("task")

func Setup() {
	go regularCleanExpiredUploads(7 * 24 * 60 * 60)
}

// regularCleanExpiredUploads 定期清理过期缓存
func regularCleanExpiredUploads(cacheSec int64) {
	ticker := time.NewTicker(1 * time.Hour)
	for {
		<-ticker.C
		removeExpiredUploads(cacheSec)
	}
	// ticker.Stop()
}

func removeExpiredUploads(cacheSec int64) {
	taskLog.Info("handle remove expired uploads")

	storageDir := strings.TrimSuffix(config.GetString("storage.dir", ""), "/")
	if storageDir == "" {
		taskLog.Error("storage dir not set")
		return
	}

	criticalTime := time.Now().Unix() - cacheSec

	// 匹配完成了，但未清理的文件
	globStr := fmt.Sprintf("%s/*/*/*/upload.file", storageDir)
	matches, err := filepath.Glob(globStr)
	if err != nil {
		taskLog.Errorf("filepath glob error: %s", err.Error())
		return
	}
	doRemoveUploads(matches, criticalTime)

	// 匹配未完成的临时文件
	globStr = fmt.Sprintf("%s/*/*/*/0.chunk", storageDir)
	matches, err = filepath.Glob(globStr)
	if err != nil {
		taskLog.Errorf("filepath glob error: %s", err.Error())
		return
	}

	doRemoveUploads(matches, criticalTime)
}

func doRemoveUploads(matches []string, criticalTime int64) {
	for _, val := range matches {
		dir := filepath.Dir(val)

		// 查看状态文件信息
		finfo, err := os.Stat(filepath.Join(dir, "stat.info"))
		if err != nil {
			taskLog.Errorw("os stat error", "error", err.Error(), "file", filepath.Join(dir, "stat.info"))
			continue
		}

		// 比较状态文件修改时间，判断是否删除临时文件
		mtime := finfo.ModTime().Unix()
		if mtime >= criticalTime {
			continue
		}

		err = os.RemoveAll(dir)
		if err != nil {
			taskLog.Errorw("os remove dir error", "error", err.Error(), "dir", dir)
		} else {
			taskLog.Infof("remove dir: %s", dir)
		}
	}
}
