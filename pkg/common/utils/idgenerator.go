package utils

import (
	"github.com/EdgarTeng/evolvest/pkg/common"
	"github.com/EdgarTeng/evolvest/pkg/common/logger"
	"os"
	"strconv"
	"sync/atomic"
)

const (
	maxVal = 1000
)

var (
	sid        int
	count      uint32
	lastMillis int64
)

func init() {
	servId := os.Getenv(common.EnvSid)
	logger.WithField(common.EnvSid, servId).Info("env")
	if servId != "" {
		if i, err := strconv.Atoi(servId); err == nil {
			sid = i
		}
	}
}

func increaseCount() uint32 {
	atomic.AddUint32(&count, 1)
	atomic.CompareAndSwapUint32(&count, maxVal, 0)
	return count
}

func GenerateId() int64 {
	millis := CurrentMillis()

	defer func() {
		lastMillis = millis
	}()
	return millis*1e6 + int64(sid*1e3) + int64(increaseCount())

}
