package utils

import (
	"math/rand"
	"time"
)

// 获取随机字符串
func GetRandString(length int) string {
	var result string
	seed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-_"
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	for i := 0; i < length; i++ {
		result += string(seed[r.Intn(len(seed) - 1)])
	}

	return result
}
