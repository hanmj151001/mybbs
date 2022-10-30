package service

import (
	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"mybbs/model"
	"mybbs/model/constants"
	"mybbs/pkg/errs"
	"mybbs/repositories"
	"mybbs/server/cache"
	"time"
)

var UserTokenService = newUserTokenService()

func newUserTokenService() *userTokenService {
	return &userTokenService{}
}

type userTokenService struct {
}

//GetUserToken 从请求体中获取UserToken
func (s *userTokenService) GetUserToken(ctx iris.Context) string {
	userToken := ctx.FormValue("userToken")
	if len(userToken) > 0 {
		return userToken
	}
	return ctx.GetHeader("X-User-Token")
}

//GetCurrent 获取当前登录用户
func (s *userTokenService) GetCurrent(ctx iris.Context) *model.User {
	token := s.GetUserToken(ctx)
	userToken := cache.UserTokenCache.Get(token)
	// 没找到授权
	if userToken == nil || userToken.Status == constants.StatusDeleted {
		return nil
	}
	// 授权过期
	if userToken.ExpiredAt <= dates.NowTimestamp() {
		return nil
	}
	user := cache.UserCache.Get(userToken.UserId)
	if user == nil || user.Status != constants.StatusOk {
		return nil
	}
	return user
}

// CheckLogin 检查登录状态
func (s *userTokenService) CheckLogin(ctx iris.Context) (*model.User, *web.CodeError) {
	user := s.GetCurrent(ctx)
	if user == nil {
		return nil, errs.NotLogin
	}
	return user, nil
}

// 生成
func (s *userTokenService) Generate(userId int64) (string, error) {
	token := strs.UUID()
	tokenExpireDays := SysConfigService.GetTokenExpireDays()
	expiredAt := time.Now().Add(time.Hour * 24 * time.Duration(tokenExpireDays))
	userToken := &model.UserToken{
		Token:      token,
		UserId:     userId,
		ExpiredAt:  dates.Timestamp(expiredAt),
		Status:     constants.StatusOk,
		CreateTime: dates.NowTimestamp(),
	}
	err := repositories.UserTokenRepository.Create(sqls.DB(), userToken)
	if err != nil {
		return "", err
	}
	return token, nil
}
