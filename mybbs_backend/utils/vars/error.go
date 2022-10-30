package vars

import (
	"collabtool/acc/utils/merror"

	"gopkg.mihoyo.com/takumi/core"
)

const (
	ZapFuncArgs = "func_args"
	ZapLogMsg   = "log_msg"
)

// 错误码定义
const (
	// 通用
	ErrParams        = -1
	ErrSystemBusy    = -1000
	ErrSystemUnknown = -1001
	ErrLogFail       = -1002
	ErrZestGetFail   = -1003

	//吕菁的err码
	ErrLanguageErr   = -10007
	ErrParamsNum     = -10008
	ErrConversionErr = -10009

	// 数据库相关
	ErrDBConn            = -10001
	ErrDBUnknownErr      = -10002
	ErrObjectNotExists   = -10003
	ErrBackupObjectOld   = -10004
	ErrDBDuplicate       = -10005
	ErrChangeCountExcept = -10006

	// 配置模块
	ErrModuleKeyEditConflict  = -20001
	ErrModuleIdNotExist       = -20002
	ErrModuleKeyNameConflict  = -20003
	ErrUpdateModuleIdNotExist = -20004
	ErrDelModuleIdNotExist    = -20005
	ErrDelKeyIdNotExist       = -20006
)

// 日志错误码定义
const (
	AlarmCacheUpdate    = 10001
	AlarmGetZestAppInfo = 10002
)

var errorMap = map[core.ErrorInfo_ErrorType]map[int]string{
	core.ErrorInfo_USER: {

		ErrConversionErr: "字符串转换失败",
		ErrLanguageErr:   "给定表达式语法无效",

		ErrParams:                "参数错误: %s",
		ErrParamsNum:             "参数数量错误",
		ErrObjectNotExists:       "操作对象(%v)不存在",
		ErrDBDuplicate:           "数据(%v)已存在",
		ErrChangeCountExcept:     "数据修改数不及预期，可能是部分数据已变更，请刷新后重试",
		ErrModuleKeyEditConflict: "配置项修改冲突，请刷新后重试嗷~",
		ErrDelKeyIdNotExist:      "删除配置项不存在",

		ErrModuleIdNotExist:       "模块不存在，需要添加模块后才能添加配置列表嗷",
		ErrModuleKeyNameConflict:  "新增feature名称冲突，请检查~",
		ErrUpdateModuleIdNotExist: "更新模块不存在",
		ErrDelModuleIdNotExist:    "删除模块不存在",
	},

	core.ErrorInfo_SYSTEM: {
		ErrSystemBusy:      "系统繁忙",
		ErrSystemUnknown:   "未知错误",
		ErrLogFail:         "日志记录失败",
		ErrDBConn:          "获取数据库连接失败",
		ErrDBUnknownErr:    "数据库未知异常",
		ErrBackupObjectOld: "备份旧版本数据失败",
		ErrZestGetFail:     "zest配置(%v)获取失败",
	},
}

func init() {
	merror.BatchRegister(errorMap)
}

// usage example
// merror.New(ctx, vars.ErrDBConn)
// merror.WithLog().WithLogCode(100).New(ctx, vars.ErrDBConn)
// merror.WithStackLog().New(ctx, vars.ErrDBConn)
