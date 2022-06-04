package gencidtask

import (
	"encoding/json"
	"fmt"
	"goframework/cmd/http/service/filesvc"

	"os"
	"path/filepath"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/jianbo-zh/go-errors"
)

var (
	ErrTimeout = errors.New("get worker timeout")
)

// 处理直接上传任务
func HandleDirectTask(mfi filesvc.MpUploadFileInfo) (res filesvc.MpUploadFileInfo, err error) {
	select {
	case ipfsCli := <-BsShells:
		// 用完了，还回去
		defer func() {
			BsShells <- ipfsCli
		}()

		res, err = handleFileTask("bs", ipfsCli, mfi)
		return

	case <-time.After(50 * time.Second):
		err = ErrTimeout
		pinLog.Errorf("[bs] wait handle file error: %s, hash: %s", err.Error(), mfi.FileHash)
		return
	}
}

// 处理队列任务
func HandleQueueTask(tag string, ipfsCli *shell.Shell, mfi filesvc.MpUploadFileInfo) {
	_, err := handleFileTask(tag, ipfsCli, mfi)
	if err != nil {
		err = filesvc.UpdatePendingHashFailed(mfi.FileHash)
		if err != nil {
			pinLog.Errorf("[%s] update failed hash error: %s, hash: %s", tag, err.Error(), mfi.FileHash)
		}
		return
	}

	err = filesvc.UpdatePendingHashSuccess(mfi.FileHash, mfi.FileCid)
	if err != nil {
		pinLog.Errorf("[%s] update success hash error: %s, hash: %s", tag, err.Error(), mfi.FileHash)
	}
}

func handleFileTask(tag string, ipfsCli *shell.Shell, mfi filesvc.MpUploadFileInfo) (
	res filesvc.MpUploadFileInfo, err error) {

	st := time.Now()

	pinLog.Infof("[%s] handle file task start, hash: %s", tag, mfi.FileHash)

	fileUploadDir, err := filesvc.GetHashUploadDir(mfi.FileHash)
	if err != nil {
		pinLog.Errorf("[%s] get upload dir error: %s, hash: %s", tag, err.Error(), mfi.FileHash)
		return
	}

	pinLog.Infof("[%s] get file upload dir: %s", tag, fileUploadDir)
	statFile := filepath.Join(fileUploadDir, "stat.info")
	statfi, err := os.OpenFile(statFile, os.O_RDWR, 0755)
	if err != nil {
		pinLog.Errorf("[%s] open stat.info error: %s, file: %s", tag, err.Error(), statFile)
		return
	}
	defer statfi.Close()
	pinLog.Infof("[%s] open file stat.info success", tag)

	fullFile := filepath.Join(fileUploadDir, "upload.file")
	originalFileCid, err := pinFileToIpfs(tag, ipfsCli, fullFile)
	if err != nil {
		pinLog.Errorf("[%s] pin upload file error: %s, file: %s", tag, err.Error(), fullFile)
		return
	}

	mfi.FileCid = originalFileCid

	mt := time.Now()

	fMime, fExt, tImgFile, cImgFile, err := getThumbnailAndCover(fullFile, fileUploadDir, mfi.FileHash)
	if err != nil {
		pinLog.Errorf("[%s] get file thumbnail and cover error: %s, file: %s", tag, err.Error(), fullFile)
		// can not return
	}
	mfi.FileMimeType = fMime
	mfi.FileExtension = fExt

	if tImgFile != "" {
		tImgCid, err := pinFileToIpfs(tag, ipfsCli, tImgFile)
		if err != nil {
			pinLog.Errorf("[%s] pin thumbnail file error: %s, file: %s", tag, err.Error(), tImgFile)
		}

		err = os.Remove(tImgFile)
		if err != nil {
			pinLog.Errorf("[%s] remove thumbnail file error: %s, file: %s", tag, err.Error(), tImgFile)
		}

		mfi.ThumbnailCid = tImgCid
	}

	if cImgFile != "" {
		cImgCid, err := pinFileToIpfs(tag, ipfsCli, cImgFile)
		if err != nil {
			pinLog.Errorf("[%s] pin cover file error: %s, file: %s", tag, err.Error(), cImgFile)
		}

		err = os.Remove(cImgFile)
		if err != nil {
			pinLog.Errorf("[%s] remove cover file error: %s, file: %s", tag, err.Error(), cImgFile)
		}

		mfi.CoverCid = cImgCid
	}

	bytes, _ := json.Marshal(mfi)
	statfi.Truncate(0)
	_, err = statfi.WriteAt(bytes, 0)
	if err != nil {
		pinLog.Errorf("[%s] update stat file error: %s, file: %s", tag, err.Error(), statFile)
		return
	}
	pinLog.Infof("[%s] update stat file success", tag)

	// 添加到ipfs之后，删除文件
	os.Remove(fullFile)

	pinLog.Infof(
		"[%s] handle file task success hash: %s, size: %d, startDt: %s, midDt: %s,"+
			" endDt: %s, midSpendTime: %f, totalSpendTime: %f",
		tag, mfi.FileHash, mfi.FileSize, st.Format(time.RFC3339), mt.Format(time.RFC3339),
		time.Now().Format(time.RFC3339), time.Since(mt).Seconds(), time.Since(st).Seconds())

	return mfi, nil
}

func pinFileToIpfs(tag string, ipfsCli *shell.Shell, file string) (cid string, err error) {

	st := time.Now()

	pinLog.Infof("[%s] pin ipfs start, file: %s", tag, file)

	osFile, err := os.Open(file)
	if err != nil {
		err = fmt.Errorf("[%s] open file error, file: %s, error: %s", tag, file, err.Error())
		return
	}
	defer osFile.Close()

	finfo, err := osFile.Stat()
	if err != nil {
		err = fmt.Errorf("[%s] stat file error, file: %s, error: %s", tag, file, err.Error())
		return
	}

	cid, err = ipfsCli.Add(osFile, shell.CidVersion(1))
	if err != nil {
		err = fmt.Errorf("[%s] pin ipfs error, file: %s, error: %s", tag, file, err.Error())
		return
	}

	pinLog.Infof("[%s] pin ipfs success, file: %s, size: %d, sec: %f", tag, file, finfo.Size(),
		time.Since(st).Seconds())

	return cid, nil
}
