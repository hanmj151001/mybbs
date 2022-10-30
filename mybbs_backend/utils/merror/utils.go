package merror

import (
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"gopkg.mihoyo.com/takumi/core"
)

// ErrorField 日志用
func ErrorField(err error) zap.Field {
	return zap.Error(FromErr(err))
}

// FromErr 解一下takumi的错误，用于一些日志输出等
func FromErr(err error) error {
	if _, ok := status.FromError(err); ok {
		return core.FromError(err)
	}

	return err
}

// CodeEq 判断code一致
func CodeEq(err error, code int32) bool {
	e := core.FromError(err)
	if e.IsRPC() {
		return e.Code == code
	}

	return false
}
