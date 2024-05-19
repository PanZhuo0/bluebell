package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

// 自定义Cliams
type MyClaims struct {
	UserID string `json:"user_id"`
	*jwt.StandardClaims
}

const timeDuration = time.Hour * 24 * 7

var mySecret = []byte("a")

func GenToken(userID string) (aToken, rToken string, err error) {
	c := MyClaims{
		UserID: userID,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeDuration).Unix(),
			Issuer:    "bluebell",
		},
	}
	// aToken
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)
	if err != nil {
		zap.L().Error("签名aToken时发生错误", zap.Error(err))
		return
	}
	// rToken
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 7 * 24).Unix(),
			Issuer:    "bluebell",
		}).SignedString(mySecret)
	if err != nil {
		zap.L().Error("签名rToken时发生错误", zap.Error(err))
		return
	}
	return
}

func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// rToken expired?
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		zap.L().Error("rToken有误", zap.Error(err))
		return
	}
	// 2.aToken expired?
	var c = new(MyClaims)
	_, err = jwt.ParseWithClaims(aToken, c, keyFunc)
	if err != nil {
		zap.L().Error("aToken:", zap.Error(err))
	}
	v, _ := err.(*jwt.ValidationError)
	if v.Errors == jwt.ValidationErrorExpired {
		zap.L().Info("aToken已过期,将尝试更新", zap.Error(err))
		return GenToken(c.UserID)
	}
	return
}

// 解析Token
func keyFunc(_ *jwt.Token) (i interface{}, err error) {
	return mySecret, nil
}

func ParseToken(tokenStr string) (claims *MyClaims, err error) {
	var token *jwt.Token
	claims = new(MyClaims)
	token, err = jwt.ParseWithClaims(tokenStr, claims, keyFunc) //使用claims 和 key 将tokenStr中的信息解析出来
	if !token.Valid {
		err = errors.New("invalid token")
	}
	return
}
