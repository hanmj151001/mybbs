package api

import (
	"github.com/kataras/iris/v12/_examples/mvc/login/services"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"
)

// 修改密码
func (c *UserController) PostUpdatePassword() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin)
	}
	var (
		oldPassword = params.FormValue(c.Ctx, "oldPassword")
		password    = params.FormValue(c.Ctx, "password")
		rePassword  = params.FormValue(c.Ctx, "rePassword")
	)
	if err := services.UserService.UpdatePassword(user.Id, oldPassword, password, rePassword); err != nil {
		return web.JsonError(err)
	}
	return web.JsonSuccess()
}
