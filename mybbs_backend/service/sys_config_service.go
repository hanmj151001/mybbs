package service

import (
	"mybbs/model/constants"
	"strconv"
)

var SysConfigService = newSysConfigService()

func newSysConfigService() *sysConfigService {
	return &sysConfigService{}
}

type sysConfigService struct {
}

func (s *sysConfigService) GetTokenExpireDays() int {
	tokenExpireDaysStr := cache.SysConfigCache.GetValue(constants.SysConfigTokenExpireDays)
	tokenExpireDays, err := strconv.Atoi(tokenExpireDaysStr)
	if err != nil {
		tokenExpireDays = constants.DefaultTokenExpireDays
	}
	if tokenExpireDays <= 0 {
		tokenExpireDays = constants.DefaultTokenExpireDays
	}
	return tokenExpireDays
}
