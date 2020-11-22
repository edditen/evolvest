package utils

import "time"

func CurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
