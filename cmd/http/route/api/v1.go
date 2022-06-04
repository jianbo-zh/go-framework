package api

import (
	"goframework/cmd/http/api/fileapi"
	"goframework/cmd/http/route/middleware"

	"github.com/gin-gonic/gin"
)

const basePath = ""

func RegisterV1(router *gin.Engine) {

	guestRouter := router.Group(basePath+"/v1", middleware.Guest)
	{
		// caRouter.POST("/user/register", userapi.Register)

		// 单文件上传接口
		guestRouter.POST("/file/sgupload", fileapi.SgUpload)
		guestRouter.POST("/file/bsgupload", fileapi.SgUploadBinary)

		// 分片上传，断点续传接口
		guestRouter.POST("/file/mpupload/init", fileapi.MpUploadInit)
		guestRouter.POST("/file/mpupload/chunk", fileapi.MpUploadChunk)
		guestRouter.POST("/file/mpupload/bchunk", fileapi.MpUploadChunkBinary)
		guestRouter.POST("/file/mpupload/finish", fileapi.MpUploadFinishV4)
		guestRouter.POST("/file/mpupload/finish2", fileapi.MpUploadFinishV4)
		guestRouter.POST("/file/mpupload/result", fileapi.MpUploadCheck)

	}
}
