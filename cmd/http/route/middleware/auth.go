package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"goframework/internal/jwt"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {

	var tokenStr string

	bareStr := ctx.GetHeader("Authorization")
	if bareStr != "" {
		bareArr := strings.Split(bareStr, " ")
		if len(bareArr) == 2 {
			tokenStr = bareArr[1]
		}
	} else {
		switch ctx.Request.Method {
		case http.MethodGet, http.MethodDelete:
			tokenStr = ctx.Request.URL.Query().Get("token")

		case http.MethodPost, http.MethodPut, http.MethodPatch:
			if strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
				if body, err := ioutil.ReadAll(ctx.Request.Body); err == nil {
					var tb map[string]interface{}
					if err := json.Unmarshal(body, &tb); err == nil {
						if _, ok := tb["token"]; ok {
							tokenStr, _ = tb["token"].(string)
						}
					}
					ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				}
			}
		}
	}

	if tokenStr == "" {
		ctx.AbortWithStatusJSON(200, map[string]interface{}{
			"code":    10000,
			"message": "token empty",
		})
		return
	}

	dataMap, err := jwt.DecodeJWT(tokenStr)
	if err != nil {
		ctx.AbortWithStatusJSON(200, map[string]interface{}{
			"code":    10000,
			"message": "parse token error",
		})
		return
	}

	id, ok := dataMap["id"]
	if !ok {
		ctx.AbortWithStatusJSON(200, map[string]interface{}{
			"code":    10000,
			"message": "token not have id",
		})
		return

	}

	userId, ok := id.(string)
	if !ok {
		ctx.AbortWithStatusJSON(200, map[string]interface{}{
			"code":    10000,
			"message": "user id type error",
		})
		return
	}

	ctx.Set("userId", userId)

	ctx.Next()
}
