package session

import (
	"fmt"
	"time"

	"github.com/K05A4B/pepper/internal/utils"
)

func getRandomSessionId() string {
	return "SESSION_ID_" + utils.GetRandString(40) + "_T" + fmt.Sprint(time.Now().Unix())
}
