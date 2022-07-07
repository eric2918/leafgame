package md5

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(str string) string {
	b := []byte(str)
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func EncodeMd5(str string) string {
	b := []byte(str)
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
