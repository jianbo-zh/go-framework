package fileapi

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"goframework/cmd/http/api/apicmn"
	"goframework/cmd/http/service/filesvc"
	"goframework/cmd/http/task/gencidtask"
	"goframework/pkg/mmux"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jianbo-zh/go-errors"
	goi18n "github.com/jianbo-zh/go-i18n"
)

// SgUpload 单文件上传
func SgUpload(ctx *gin.Context) {
	var params SgUploadParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		apicmn.Error(ctx, errors.Newc(apicmn.ErrParamErr, goi18n.Sprintf("params error")).With(errors.Inner(err)))
		return
	}

	apiLogger.Infof("[bs] sgupload start, hash: %s", params.Hash)

	stime := time.Now()

	file, err := ctx.FormFile("file")
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("get form file error")).With(errors.Inner(err)))
		return
	}

	ofile, err := file.Open()
	if err != nil {
		apicmn.Error(ctx, errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("file open error")).With(errors.Inner(err)))
		return
	}
	defer ofile.Close()

	ufbs, err := io.ReadAll(ofile)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("read upload file error")).With(errors.Inner(err)))
		return
	}

	if params.Hash != fmt.Sprintf("%x", sha1.Sum(ufbs)) {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadFinish, goi18n.Sprintf("upload file hash error: %s", params.Hash)))
		return
	}

	fileUploadDir, err := filesvc.GetHashUploadDir(params.Hash)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("get upload dir error")).With(errors.Inner(err)))
		return
	}

	fi, err := os.Stat(fileUploadDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(fileUploadDir, 0755)
			if err != nil {
				apicmn.Error(ctx,
					errors.Newc(apicmn.ErrUploadInit,
						goi18n.Sprintf("make file save dir error")).With(errors.Inner(err)))
				return
			}
		}
	} else if !fi.IsDir() {
		apicmn.Error(ctx, errors.Newc(apicmn.ErrUploadInit, goi18n.Sprintf("file save dir is not dir error")))
		return
	}

	apiLogger.Infof("[bs] sgupload prepared, hash: %s", params.Hash)

	// TODO: 互斥锁
	mmux.Lock(params.Hash)
	defer mmux.Unlock(params.Hash)

	apiLogger.Infof("[bs] sgupload get lock, hash: %s", params.Hash)

	statfi, err := os.OpenFile(filepath.Join(fileUploadDir, "stat.info"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("open file save dir error")).With(errors.Inner(err)))
		return
	}
	defer statfi.Close()

	bytes, err := io.ReadAll(statfi)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("read file info error")).With(errors.Inner(err)))
		return
	}

	var mfi filesvc.MpUploadFileInfo
	if len(bytes) == 0 {
		// new
		mfi.FileName = params.FileName
		mfi.FileSize = params.FileSize
		mfi.FileHash = params.Hash
		mfi.ChunkSize = ChunkSize
		mfi.UploadChunks = make([]filesvc.MpUploadChunkInfo, 0)

		bytes, _ = json.Marshal(mfi)
		_, err = statfi.Write(bytes)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadInit, goi18n.Sprintf("save file info error")).With(errors.Inner(err)))
			return
		}

	} else {
		err = json.Unmarshal(bytes, &mfi)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadInit,
					goi18n.Sprintf("json unmarshal file info error")).With(errors.Inner(err)))
			return
		}

		if mfi.FileName != params.FileName || mfi.FileSize != params.FileSize {
			mfi.FileName = params.FileName
			mfi.FileSize = params.FileSize

			bytes, _ = json.Marshal(mfi)
			statfi.Truncate(0)
			_, err = statfi.WriteAt(bytes, 0)
			if err != nil {
				apicmn.Error(ctx,
					errors.Newc(apicmn.ErrUploadInit, goi18n.Sprintf("save file info error")).With(errors.Inner(err)))
				return
			}
		}
	}

	apiLogger.Infof("[bs] sgupload update stat, hash: %s", params.Hash)

	if mfi.FileCid == "" {

		fullFile := filepath.Join(fileUploadDir, "upload.file")

		err = ctx.SaveUploadedFile(file, fullFile)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("save upload file error")).With(errors.Inner(err)))
			return
		}

		apiLogger.Infof("[bs] sgupload save upload file, hash: %s", params.Hash)

		mfi, err = gencidtask.HandleDirectTask(mfi)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("save upload file error")).With(errors.Inner(err)))
			return
		}

		pinLog.Infof(
			"sgUpload success hash: %s, size: %d, spendTime: %f",
			mfi.FileHash, mfi.FileSize, time.Since(stime).Seconds())
	}

	apicmn.Success(ctx, mfi)
}

// SgUploadBinary 单文件上传
func SgUploadBinary(ctx *gin.Context) {
	var params SgUploadParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		apicmn.Error(ctx, errors.Newc(apicmn.ErrParamErr, goi18n.Sprintf("params error")).With(errors.Inner(err)))
		return
	}

	stime := time.Now()

	ufbs, err := ctx.GetRawData()
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("read upload file error")).With(errors.Inner(err)))
		return
	}

	if params.Hash != fmt.Sprintf("%x", sha1.Sum(ufbs)) {
		apicmn.Error(
			ctx,
			errors.Newc(
				apicmn.ErrUploadFinish,
				goi18n.Sprintf("upload file hash error: %s", params.Hash)).With(errors.Inner(err)))
		return
	}

	fileUploadDir, err := filesvc.GetHashUploadDir(params.Hash)
	if err != nil {
		apicmn.Error(
			ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("get upload dir error")).With(errors.Inner(err)))
		return
	}

	fi, err := os.Stat(fileUploadDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(fileUploadDir, 0755)
			if err != nil {
				apicmn.Error(ctx,
					errors.Newc(apicmn.ErrUploadInit,
						goi18n.Sprintf("make file save dir error")).With(errors.Inner(err)))
				return
			}
		}
	} else if !fi.IsDir() {
		apicmn.Error(ctx, errors.Newc(apicmn.ErrUploadInit, goi18n.Sprintf("file save dir is not dir error")))
		return
	}

	// TODO: 互斥锁
	mmux.Lock(params.Hash)
	defer mmux.Unlock(params.Hash)

	statfi, err := os.OpenFile(filepath.Join(fileUploadDir, "stat.info"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("open file save dir error")).With(errors.Inner(err)))
		return
	}
	defer statfi.Close()

	bytes, err := io.ReadAll(statfi)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("read file info error")).With(errors.Inner(err)))
		return
	}

	var mfi filesvc.MpUploadFileInfo
	if len(bytes) == 0 {
		// new
		mfi.FileName = params.FileName
		mfi.FileSize = params.FileSize
		mfi.FileHash = params.Hash
		mfi.ChunkSize = ChunkSize
		mfi.UploadChunks = make([]filesvc.MpUploadChunkInfo, 0)

		bytes, _ = json.Marshal(mfi)
		_, err = statfi.Write(bytes)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadInit, goi18n.Sprintf("save file info error")).With(errors.Inner(err)))
			return
		}

	} else {
		err = json.Unmarshal(bytes, &mfi)
		if err != nil {
			apicmn.Error(
				ctx,
				errors.Newc(apicmn.ErrUploadInit,
					goi18n.Sprintf("json unmarshal file info error")).With(errors.Inner(err)))
			return
		}

		if mfi.FileName != params.FileName || mfi.FileSize != params.FileSize {
			mfi.FileName = params.FileName
			mfi.FileSize = params.FileSize

			bytes, _ = json.Marshal(mfi)
			statfi.Truncate(0)
			_, err = statfi.WriteAt(bytes, 0)
			if err != nil {
				apicmn.Error(ctx,
					errors.Newc(apicmn.ErrUploadInit, goi18n.Sprintf("save file info error")).With(errors.Inner(err)))
				return
			}
		}
	}

	if mfi.FileCid == "" {

		fullFile := filepath.Join(fileUploadDir, "upload.file")

		err = os.WriteFile(fullFile, ufbs, 0755)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("save upload file error")).With(errors.Inner(err)))
			return
		}

		mfi, err = gencidtask.HandleDirectTask(mfi)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("save upload file error")).With(errors.Inner(err)))
			return
		}

		pinLog.Infof(
			"sgUploadB success cid: %s, size: %d, spendTime: %f",
			mfi.FileCid, mfi.FileSize, time.Since(stime).Seconds())
	}

	apicmn.Success(ctx, mfi)
}
