package merror

import (
	"collabtool/acc/utils/i18n"
	"context"
	"fmt"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"gopkg.mihoyo.com/takumi/core"
	"gopkg.mihoyo.com/takumi/log"
	"gopkg.mihoyo.com/takumi/settings/zest"
)

var (
	mErrors map[int]mError // 注册的错误字典
	mu      sync.Mutex     // 互斥锁 理论上初始信息只在初始化的时候注入，但是防止并发注册的情况
	once    sync.Once      // 配置初始化限制
	mConfig *config        // 配置
)

const (
	// ErrInternalSystemCode 内部错误码定义
	ErrInternalSystemCode = -500
	// ErrInternalSystemMsg 未知错误key
	ErrInternalSystemMsg = "unknown error code: %d"
)

// MError 错误封装接口
type MError interface {
	New(ctx context.Context, errorCode int, args ...interface{}) error
	WithLog(fields ...zap.Field) MError
	WithLogCode(code int) MError
	WithStackSkip(skip int) MError
	WithCallSkip(skip int) MError
	WithStackLog(fields ...zap.Field) MError
}

type config struct {
	EnableI18n     bool   `yaml:"enable_i18n"`     // 启用多语言
	Mi18nNamespace string `yaml:"mi18n_namespace"` // 多语言key的命名空间
}

// 错误基础数据结构
type mError struct {
	errorType     core.ErrorInfo_ErrorType
	code          int
	defaultFormat string
}

// 构建错误的参数结构
type errorCore struct {
	errorCode int
	args      []interface{}
	logMeta   *logMeta
	callSkip  int
	stackSkip int
}

// 打印日志所需参数结构
type logMeta struct {
	logCode    int
	logField   []log.Field
	isLogStack bool
}

func init() {
	mErrors = make(map[int]mError, 0)
	mu = sync.Mutex{}
	once = sync.Once{}
}

// Register 注册错误信息，理论上只在启动初始化阶段的时候进行调用
func Register(errorCode int, errorType core.ErrorInfo_ErrorType, defaultFormat string) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := mErrors[errorCode]; ok {
		panic(fmt.Sprintf("duplicate registration error code: %d", errorCode))
	}

	mErrors[errorCode] = mError{
		errorType:     errorType,
		code:          errorCode,
		defaultFormat: defaultFormat,
	}
}

// BatchRegister 批量注册
func BatchRegister(errorMap map[core.ErrorInfo_ErrorType]map[int]string) {
	mu.Lock()
	defer mu.Unlock()

	for errorType, codeMap := range errorMap {
		for code, defaultFormat := range codeMap {
			if _, ok := mErrors[code]; ok {
				panic(fmt.Sprintf("duplicate registration error code: %d", code))
			}

			mErrors[code] = mError{
				errorType:     errorType,
				code:          code,
				defaultFormat: defaultFormat,
			}
		}
	}
}

// New 创建错误
func New(ctx context.Context, errorCode int, args ...interface{}) error {
	return newCore().New(ctx, errorCode, args...)
}

// WithLog 在创建错误的时候打印日志
func WithLog(fields ...zap.Field) MError {
	return newCore().WithLog(fields...)
}

// WithStackLog 在创建错误的时候打印堆栈日志
func WithStackLog(fields ...zap.Field) MError {
	return newCore().WithStackLog(fields...)
}

// New 创建错误
func (e *errorCore) New(ctx context.Context, errorCode int, args ...interface{}) error {
	return e.withArgs(1, errorCode, args...).build(ctx)
}

// WithLog 打印日志
func (e *errorCore) WithLog(fields ...zap.Field) MError {
	cloned := *e
	if cloned.logMeta == nil {
		cloned.logMeta = newDefaultLogMeta()
	}

	cloned.logMeta.logField = fields
	return &cloned
}

// WithLogCode 设置日志Code
func (e *errorCore) WithLogCode(code int) MError {
	cloned := *e
	if cloned.logMeta == nil {
		cloned.logMeta = newDefaultLogMeta()
	}

	cloned.logMeta.logCode = code
	return &cloned
}

// WithStackSkip 设置日志堆栈Skip
func (e *errorCore) WithStackSkip(skip int) MError {
	cloned := *e
	cloned.stackSkip = skip
	return &cloned
}

// WithCallSkip 设置日志CallSkip
func (e *errorCore) WithCallSkip(skip int) MError {
	cloned := *e
	cloned.callSkip = skip
	return &cloned
}

// WithStackLog 打印带堆栈信息的日志
func (e *errorCore) WithStackLog(fields ...zap.Field) MError {
	cloned := *e
	if cloned.logMeta == nil {
		cloned.logMeta = newDefaultLogMeta()
	}

	cloned.logMeta.isLogStack = true
	cloned.logMeta.logField = fields
	return &cloned
}

// withArgs 注入参数
func (e *errorCore) withArgs(callSkip int, errorCode int, args ...interface{}) *errorCore {
	ne := *e
	ne.callSkip += callSkip
	ne.stackSkip += callSkip
	ne.errorCode = errorCode
	ne.args = args
	return &ne
}

// build 构建错误
func (e *errorCore) build(ctx context.Context) error {
	initLazy()
	meta, ok := mErrors[e.errorCode]
	if !ok {
		return core.SystemError(ErrInternalSystemCode, fmt.Sprintf(ErrInternalSystemMsg, e.errorCode))
	}

	var msg string
	// 多语言支持
	if mConfig.EnableI18n {
		msg = i18n.Get(mConfig.Mi18nNamespace).T(ctx, strconv.Itoa(e.errorCode), meta.defaultFormat, e.args...)
	} else {
		msg = fmt.Sprintf(meta.defaultFormat, e.args...)
	}

	// 日志打印
	if e.logMeta != nil {
		e.callSkip++
		e.stackSkip++
		e.printLog(ctx, msg, e.args...)
	}

	switch meta.errorType {
	case core.ErrorInfo_SYSTEM:
		return core.SystemError(meta.code, msg)
	case core.ErrorInfo_USER:
		return core.UserError(meta.code, msg)
	case core.ErrorInfo_UNKNOWN:
		return core.UnknownError(meta.code, msg)
	}

	return core.UnknownError(meta.code, msg)
}

// printLog 打印错误日志
func (e *errorCore) printLog(ctx context.Context, retMsg string, args ...interface{}) {
	initLazy()
	meta, ok := mErrors[e.errorCode]
	if !ok {
		return
	}

	msg := fmt.Sprintf(meta.defaultFormat, args...)

	fields := append([]log.Field{log.String("retMsg", retMsg)}, e.logMeta.logField...)

	logIns := log.Get().CallSkip(e.callSkip + 1)
	if e.logMeta.logCode != int(log.CodeErr) {
		logIns = logIns.Code(e.logMeta.logCode)
	} else {
		logIns = logIns.Code(e.errorCode) // 默认使用错误Code
	}
	if e.logMeta.isLogStack {
		logIns.StackSkip(e.stackSkip+1).Stack(ctx, msg, fields...)
	} else {
		logIns.Error(ctx, msg, fields...)
	}
}

// newCore 创建
func newCore() *errorCore {
	return &errorCore{}
}

// newDefaultLogMeta 创建默认logMeta
func newDefaultLogMeta() *logMeta {
	return &logMeta{
		logCode:    int(log.CodeErr),
		logField:   make([]log.Field, 0),
		isLogStack: false,
	}
}

// initLazy 配置懒加载
func initLazy() {
	once.Do(func() {
		defer func() {
			if e := recover(); e != nil {
				mConfig = &config{}
				return
			}
		}()

		c := config{}
		zest.Get("m_error", &c)
		mConfig = &c
	})
}
