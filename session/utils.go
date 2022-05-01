package session

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/K05A4B/pepper/internal/utils"
)

func getRandomSessionId() string {
	timeString := fmt.Sprint(time.Now().Unix())
	timeStringRune := []rune(timeString)

	timeFeild := ""

	for i := 5; i > 0; i-- {
		char := timeStringRune[len(timeStringRune)-i]
		timeFeild += string(char)
	}

	return "SESSION_ID_" + utils.GetRandString(40) + "_T" + timeFeild
}

// 转成 session 对象
func binaryToSession(b []byte) (*Session, error) {
	sess := newSession()

	reader := bytes.NewReader(b)

	decoder := gob.NewDecoder(reader)
	return sess, decoder.Decode(&sess.Value)
}