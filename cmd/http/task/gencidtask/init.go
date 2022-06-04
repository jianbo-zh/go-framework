package gencidtask

import (
	"goframework/cmd/http/service/filesvc"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/jianbo-zh/go-log"
)

// 任务通道
var finishChannel = make(chan filesvc.MpUploadFileInfo, 1000)
var largeFileChannel = make(chan filesvc.MpUploadFileInfo, 1000)
var normalFileChannel = make(chan filesvc.MpUploadFileInfo, 1000)
var smallFileChannel = make(chan filesvc.MpUploadFileInfo, 1000)

// 处理shell
var alShells []*shell.Shell //大文件处理 ipfs client
var anShells []*shell.Shell //普通文件处理 ipfs client
var asShells []*shell.Shell //小文件处理 ipfs client

var BsShells chan *shell.Shell // 小文件通道

// 日志记录器
var pinLog = log.Logger("pins")
var checkLog = log.Logger("check")
var imgLog = log.Logger("image")
