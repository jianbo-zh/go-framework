package apicmn

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jianbo-zh/go-errors"
)

// ResponseModel reponse container
type Response struct {
	Code         int         `json:"code"`
	Message      string      `json:"message"`
	SubMessage   string      `json:"submessage"`
	Data         interface{} `json:"data"`
	ResponseTime int64       `json:"rsptime"`
}

func Success(ctx *gin.Context, datas ...interface{}) {
	ctx.JSON(200, buildSuccess(datas...))
}

func Error(ctx *gin.Context, err error) {
	ctx.JSON(200, buildError(err))
}

// buildSuccess 构造成功响应
func buildSuccess(datas ...interface{}) Response {
	return Response{
		Code:         200,
		Message:      "ok",
		ResponseTime: time.Now().Unix(),
	}
}

// buildError 构造失败响应
func buildError(err error) Response {

	var msg, submsg string
	var code int

	if err1, ok := err.(*errors.Error); ok {
		code = err1.Code
		msg = err1.Messge
		submsg = err1.Error()

	} else {
		code = 10000
		msg = err.Error()
		submsg = msg
	}

	return Response{
		Code:         code,
		Message:      msg,
		SubMessage:   submsg,
		ResponseTime: time.Now().Unix(),
	}
}
