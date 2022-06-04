package fileapi

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"goframework/cmd/http/api/apicmn"
	"goframework/cmd/http/service/filesvc"
	"goframework/pkg/mmux"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/gin-gonic/gin"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/jianbo-zh/go-errors"
	goi18n "github.com/jianbo-zh/go-i18n"
)

const (
	ChunkSize = 4194304 // 4MB
)

var incrNum int64 = 0

var shells []*shell.Shell

// MpUploadInit 初始化上传
func MpUploadInit(ctx *gin.Context) {
	var params MpUploadInitParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrParamErr, goi18n.Sprintf("params error")).With(errors.Inner(err)))
		return
	}

	fileUploadDir, err := filesvc.GetHashUploadDir(params.Hash)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadInit,
				goi18n.Sprintf("get upload error")).With(errors.Inner(err)))
		return
	}

	fi, err := os.Stat(fileUploadDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(fileUploadDir, 0755)
			if err != nil {
				apicmn.Error(
					ctx,
					errors.Newc(
						apicmn.ErrUploadInit,
						goi18n.Sprintf("make file save dir error")).With(errors.Inner(err)),
				)
				return
			}
		}
	} else if !fi.IsDir() {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadInit,
				goi18n.Sprintf("file save dir is not dir error")).With(errors.Inner(err)))
		return
	}

	// key互斥锁
	mmux.Lock(params.Hash)
	defer mmux.Unlock(params.Hash)

	f, err := os.OpenFile(filepath.Join(fileUploadDir, "stat.info"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadInit,
				goi18n.Sprintf("open file save dir error")).With(errors.Inner(err)))
		return
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadInit,
				goi18n.Sprintf("read file info error")).With(errors.Inner(err)))
		return
	}

	var mfi filesvc.MpUploadFileInfo
	if len(bytes) == 0 {
		// new
		mfi.FileName = params.FileName
		mfi.FileSize = params.FileSize
		mfi.ChunkSize = ChunkSize
		mfi.FileHash = params.Hash
		mfi.UploadState = 0
		mfi.UploadChunks = make([]filesvc.MpUploadChunkInfo, 0)

		bytes, _ = json.Marshal(mfi)
		_, err = f.Write(bytes)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadInit,
					goi18n.Sprintf("save file info error")).With(errors.Inner(err)))
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
			f.Truncate(0)
			_, err = f.WriteAt(bytes, 0)
			if err != nil {
				apicmn.Error(ctx,
					errors.Newc(apicmn.ErrUploadInit,
						goi18n.Sprintf("save file info error")).With(errors.Inner(err)))
				return
			}
		}
	}

	apicmn.Success(ctx, mfi)
}

// MpUploadChunk 上传分片
func MpUploadChunk(ctx *gin.Context) {
	var params MpUploadChunkParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrParamErr, goi18n.Sprintf("params error")).With(errors.Inner(err)))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("get form file error")).With(errors.Inner(err)))
		return
	}

	fileUploadDir, err := filesvc.GetHashUploadDir(params.Hash)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("get upload dir error")).With(errors.Inner(err)))
		return
	}

	// TODO: 互斥锁
	mmux.Lock(params.Hash)
	defer mmux.Unlock(params.Hash)

	f, err := os.OpenFile(filepath.Join(fileUploadDir, "stat.info"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("open file save dir error")).With(errors.Inner(err)))
		return
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("read file info error")).With(errors.Inner(err)))
		return
	}

	var mfi filesvc.MpUploadFileInfo
	err = json.Unmarshal(bytes, &mfi)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("json unmarshal file info error")).With(errors.Inner(err)))
		return
	}

	// 检查是否已经上传完成
	if mfi.FileCid != "" || mfi.UploadState == 1 {
		// 已经添加到ipfs或上传完成了，则不处理，直接返回结果
		apicmn.Success(ctx, mfi)
		return
	}

	chunkFilePath := filepath.Join(fileUploadDir, fmt.Sprintf("%d.chunk", params.ChunkIndex))

	err = ctx.SaveUploadedFile(file, chunkFilePath)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("save upload file error")).With(errors.Inner(err)))
		return
	}

	// 检查分片是否已经上传
	var chunkExisted bool
	for i, v := range mfi.UploadChunks {
		if v.ChunkIndex == params.ChunkIndex {
			mfi.UploadChunks[i].ChunkSize = params.ChunkSize
			chunkExisted = true
			break
		}
	}

	if !chunkExisted {
		mfi.UploadChunks = append(mfi.UploadChunks, filesvc.MpUploadChunkInfo{
			ChunkIndex: params.ChunkIndex,
			ChunkSize:  params.ChunkSize,
		})

		sort.Slice(mfi.UploadChunks, func(i, j int) bool {
			return mfi.UploadChunks[i].ChunkIndex < mfi.UploadChunks[j].ChunkIndex
		})
	}

	bytes, _ = json.Marshal(mfi)
	f.Truncate(0)
	_, err = f.WriteAt(bytes, 0)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("save file info error")).With(errors.Inner(err)))
		return
	}

	apicmn.Success(ctx, mfi)
}

// MpUploadChunk 上传分片
func MpUploadChunkBinary(ctx *gin.Context) {
	var params MpUploadChunkParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrParamErr, goi18n.Sprintf("params error")).With(errors.Inner(err)))
		return
	}

	file, err := ctx.GetRawData()
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("get form file error")).With(errors.Inner(err)))
		return
	}

	fileUploadDir, err := filesvc.GetHashUploadDir(params.Hash)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("get upload dir error")).With(errors.Inner(err)))
		return
	}

	// TODO: 互斥锁
	mmux.Lock(params.Hash)
	defer mmux.Unlock(params.Hash)

	f, err := os.OpenFile(filepath.Join(fileUploadDir, "stat.info"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("open file save dir error")).With(errors.Inner(err)))
		return
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("read file info error")).With(errors.Inner(err)))
		return
	}

	var mfi filesvc.MpUploadFileInfo
	err = json.Unmarshal(bytes, &mfi)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk,
				goi18n.Sprintf("json unmarshal file info error")).With(errors.Inner(err)))
		return
	}

	// 检查是否已经上传完成
	if mfi.FileCid != "" || mfi.UploadState == 1 {
		// 已经添加到ipfs或上传完成了，则不处理，直接返回结果
		apicmn.Success(ctx, mfi)
		return
	}

	// 保存上传文件
	chunkFilePath := filepath.Join(fileUploadDir, fmt.Sprintf("%d.chunk", params.ChunkIndex))
	err = os.WriteFile(chunkFilePath, file, 0755)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("save upload file error")).With(errors.Inner(err)))
		return
	}

	// 检查分片是否已经上传
	var chunkExisted bool
	for i, v := range mfi.UploadChunks {
		if v.ChunkIndex == params.ChunkIndex {
			mfi.UploadChunks[i].ChunkSize = params.ChunkSize
			chunkExisted = true
			break
		}
	}

	if !chunkExisted {
		mfi.UploadChunks = append(mfi.UploadChunks, filesvc.MpUploadChunkInfo{
			ChunkIndex: params.ChunkIndex,
			ChunkSize:  params.ChunkSize,
		})

		sort.Slice(mfi.UploadChunks, func(i, j int) bool {
			return mfi.UploadChunks[i].ChunkIndex < mfi.UploadChunks[j].ChunkIndex
		})
	}

	bytes, _ = json.Marshal(mfi)
	f.Truncate(0)
	_, err = f.WriteAt(bytes, 0)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadChunk, goi18n.Sprintf("save file info error")).With(errors.Inner(err)))
		return
	}

	apicmn.Success(ctx, mfi)
}

// MpUploadFinishV4 上传完成，加入处理队列
func MpUploadFinishV4(ctx *gin.Context) {
	var params MpUploadFinishParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		apicmn.Error(ctx, errors.Newc(apicmn.ErrParamErr, goi18n.Sprintf("params error")).With(errors.Inner(err)))
		return
	}

	fileUploadDir, err := filesvc.GetHashUploadDir(params.Hash)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadFinish, goi18n.Sprintf("get upload dir error")).With(errors.Inner(err)))
		return
	}

	// TODO: 互斥锁
	mmux.Lock(params.Hash)
	defer mmux.Unlock(params.Hash)

	fi, err := os.OpenFile(filepath.Join(fileUploadDir, "stat.info"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadFinish, goi18n.Sprintf("open file save dir error")).With(errors.Inner(err)))
		return
	}
	defer fi.Close()

	bytes, err := io.ReadAll(fi)
	if err != nil {
		apicmn.Error(ctx,
			errors.Newc(apicmn.ErrUploadFinish, goi18n.Sprintf("read file info error")).With(errors.Inner(err)))
		return
	}

	var mfi filesvc.MpUploadFileInfo
	err = json.Unmarshal(bytes, &mfi)
	if err != nil {
		apicmn.Error(
			ctx,
			errors.Newc(apicmn.ErrUploadFinish,
				goi18n.Sprintf("json unmarshal file info error")).With(errors.Inner(err)),
		)
		return
	}

	if mfi.UploadState == 0 {
		// 合并
		fullFile := filepath.Join(fileUploadDir, "upload.file")
		ff, err := os.OpenFile(fullFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755)
		if err != nil {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadFinish, goi18n.Sprintf("open full file error")).With(errors.Inner(err)))
			return
		}
		defer ff.Close()

		hashSum := sha1.New()

		multiWriter := io.MultiWriter(ff, hashSum)

		for idx, chunk := range mfi.UploadChunks {
			if idx != int(chunk.ChunkIndex) {
				apicmn.Error(ctx,
					errors.Newc(apicmn.ErrUploadFinish,
						goi18n.Sprintf("chunk file [%d] missing", idx)).With(errors.Inner(err)))
				return
			}

			chunkFile := filepath.Join(fileUploadDir, fmt.Sprintf("%d.chunk", chunk.ChunkIndex))
			fc, err := os.OpenFile(chunkFile, os.O_RDONLY, 0755)
			if err != nil {
				apicmn.Error(ctx,
					errors.Newc(apicmn.ErrUploadFinish,
						goi18n.Sprintf("open chunk file [%d] error", idx)).With(errors.Inner(err)))
				return
			}

			_, err = io.Copy(multiWriter, fc)
			fc.Close() // after copy, do close

			if err != nil {
				apicmn.Error(ctx,
					errors.Newc(apicmn.ErrUploadFinish,
						goi18n.Sprintf("copy chunk file [%d] error", idx)).With(errors.Inner(err)))
				return
			}
		}

		// 校验 hash是否正确
		if params.Hash != fmt.Sprintf("%x", hashSum.Sum(nil)) {
			apicmn.Error(ctx,
				errors.Newc(apicmn.ErrUploadFinish,
					goi18n.Sprintf("upload file hash error: %s", params.Hash)).With(errors.Inner(err)))
			return
		}

		// 文件合并完成，更新状态
		oldUploadChunks := mfi.UploadChunks

		mfi.UploadState = 1
		mfi.FileHash = params.Hash
		mfi.UploadChunks = make([]filesvc.MpUploadChunkInfo, 0)

		bytes, _ = json.Marshal(mfi)
		fi.Truncate(0)
		_, err = fi.WriteAt(bytes, 0)
		if err != nil {
			apicmn.Error(
				ctx,
				errors.Newc(apicmn.ErrUploadFinish, goi18n.Sprintf("update stat file error")).With(errors.Inner(err)),
			)
			return
		}

		// 删除所有分片
		for _, chunk := range oldUploadChunks {
			os.Remove(filepath.Join(fileUploadDir, fmt.Sprintf("%d.chunk", chunk.ChunkIndex)))
		}
	}

	if mfi.FileCid != "" {
		apicmn.Success(ctx, mfi)
		return
	}

	err = filesvc.AddPendHash(mfi.FileHash, mfi.FileSize)
	if err != nil {
		apicmn.Error(
			ctx,
			errors.Newc(apicmn.ErrUploadFinish, goi18n.Sprintf("push file to db queue failed")).With(errors.Inner(err)))
		return
	}

	apicmn.Success(ctx, mfi)
}

// MpUploadCheck 检查上传结果
func MpUploadCheck(ctx *gin.Context) {
	var params MpUploadCheckParams
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		apicmn.Error(ctx, errors.Newc(apicmn.ErrParamErr, goi18n.Sprintf("params error")).With(errors.Inner(err)))
		return

	} else if len(params.Hashs) == 0 {
		apicmn.Error(ctx, errors.Newc(apicmn.ErrParamErr, "params error2"))
		return
	}

	resp := make([]MpUploadCheckInfo, len(params.Hashs))

	for i, hash := range params.Hashs {
		mfr := MpUploadCheckInfo{
			FileHash: hash,
		}
		mfi, err := filesvc.GetHashMfi(hash)
		if err != nil {
			mfr.Error = err.Error()

		} else {
			mfr.FileName = mfi.FileName
			mfr.FileSize = mfi.FileSize
			mfr.FileCid = mfi.FileCid
			mfr.ThumbnailCid = mfi.ThumbnailCid
			mfr.CoverCid = mfi.CoverCid
			mfr.FileMimeType = mfi.FileMimeType
			mfr.FileExtension = mfi.FileExtension
			mfr.FileHash = mfi.FileHash
		}

		resp[i] = mfr
	}

	apicmn.Success(ctx, resp)
}
