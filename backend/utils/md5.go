package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(str string) (ret string) {
	hash := md5.New()
	hash.Write(mySecret)
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
