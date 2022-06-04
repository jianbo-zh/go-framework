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

func Guest(ctx *gin.Context) {

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
		// passthrogh
		ctx.Next()
		return
	}

	dataMap, err := jwt.DecodeJWT(tokenStr)
	if err == nil {
		if id, ok := dataMap["id"]; ok {
			if userId, ok := id.(string); ok {
				ctx.Set("userId", userId)
			}
		}
	}

	ctx.Next()
}
