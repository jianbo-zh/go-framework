package jwt

import (
	"time"

	"github.com/jianbo-zh/go-errors"

	"github.com/dgrijalva/jwt-go"
)

type LoginClaims struct {
	Id string `json:"id"`

	jwt.StandardClaims
}

const (
	jwtSigKey = "adsfa#^$%#$fgrf"
)

func EncodeJWT(data map[string]interface{}, expireDuration int64) (string, error) {

	calim := jwt.MapClaims(data)

	if expireDuration < 0 {
		calim["exp"] = time.Now().Unix() + expireDuration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, calim)

	sigStr, err := token.SignedString(jwtSigKey)
	if err != nil {
		return "", errors.New("token signed string error").With(errors.Inner(err))
	}

	return sigStr, nil
}

func DecodeJWT(tokenStr string) (map[string]interface{}, error) {

	jwtToken, err := jwt.Parse(tokenStr, func(*jwt.Token) (interface{}, error) {
		return jwtSigKey, nil
	})
	if err != nil {
		return nil, errors.New("parse token error").With(errors.Inner(err))
	}

	if !jwtToken.Valid {
		return nil, errors.New("token invalid")
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); !ok {
		return nil, errors.New("claims is not mapclaims")

	} else {
		return claims, nil
	}
}
