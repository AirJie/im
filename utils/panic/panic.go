package panic

import (
	"github.com/aijie/michat/server/logger"
	"go.uber.org/zap"
)

func RecoverPanic() {
	err := recover()
	if err != nil {
		logger.Logger.DPanic("panic", zap.Any("panic", err), zap.Stack("stack"))
	}
}
//
//func GetStackInfo() string {
//	buf := make([]byte, 4095)
//	n := runtime.Stack(buf, false)
//	return string(buf[0:n])
//}
