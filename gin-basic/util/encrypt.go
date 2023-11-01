package util

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func Md5String(s string) string {
	timeStr := string(time.Now().UnixNano())
	timeStr += s
	ctx := md5.New()
	ctx.Write([]byte(s))
	return hex.EncodeToString(ctx.Sum(nil))
}
