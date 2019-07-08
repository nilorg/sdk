package random

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
)

// SmsCode 短信验证码
func SmsCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	var buffer bytes.Buffer
	for i := 0; i < length; i++ {
		buffer.WriteString(strconv.Itoa(rand.Intn(10)))
	}
	return buffer.String()
}
