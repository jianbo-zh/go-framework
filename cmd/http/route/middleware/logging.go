package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"goframework/pkg/snow"

	"github.com/gin-gonic/gin"
	"github.com/jianbo-zh/go-log"
)

var routeLogger = log.Logger("middlewares")

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Log(ctx *gin.Context) {

	startTime := time.Now()

	data, err := ctx.GetRawData()
	if err != nil {
		routeLogger.Errorw(
			fmt.Sprintf("request uri: %s", ctx.Request.RequestURI),
			"method", ctx.Request.Method,
			"clientIP", ctx.ClientIP(),
			"userAgent", ctx.Request.UserAgent(),
			"timestamp", startTime.Format(time.RFC3339),
			"error", fmt.Errorf("get raw data error: %w", err),
		)

		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//读过的字节流重新放到 body
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	var requestBody string
	if len(data) > 200 {
		requestBody = string(data[0:200]) + "..."
	} else {
		requestBody = string(data)
	}

	requestID := snow.SnowId()

	// 请求日志
	routeLogger.Infow(
		fmt.Sprintf("request uri: %s", ctx.Request.RequestURI),
		"requestID", requestID,
		"method", ctx.Request.Method,
		"clientIP", ctx.ClientIP(),
		"userAgent", ctx.Request.UserAgent(),
		"timestamp", startTime.Format(time.RFC3339),
		"requestBody", requestBody,
	)

	bodyLogWriter := &bodyLogWriter{
		body:           bytes.NewBufferString(""),
		ResponseWriter: ctx.Writer,
	}
	ctx.Writer = bodyLogWriter

	//处理请求
	ctx.Next()

	//响应日志
	responseBody := bodyLogWriter.body.String()
	if len(responseBody) > 200 {
		responseBody = responseBody[0:200] + "..."
	}

	routeLogger.Infow(
		fmt.Sprintf("response uri: %s", ctx.Request.RequestURI),
		"requestID", requestID,
		"responseBody", responseBody,
		"responseSize", ctx.Writer.Size(),
		"statusCode", ctx.Writer.Status(),
		"error", ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
		"timestamp", time.Now().Format(time.RFC3339),
		"delayTime", fmt.Sprintf("%vms", time.Since(startTime).Milliseconds()),
	)
}
