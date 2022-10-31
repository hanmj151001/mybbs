package logic

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/models"
	"bluebell_backend/pkg/gofunc"
	"bluebell_backend/pkg/jwt"
	"bluebell_backend/pkg/snowflake"
	"bluebell_backend/pkg/valiadate"
	"database/sql"
	"errors"
	"strings"
	//"go.uber.org/zap"
)

/**
 * @Author mengjie.han
 * @Description //TODO 存放注册业务逻辑的代码
 * @Date 21:52 2022/2/10
 **/
//SignUp

////TODO 这里返回一个唯一的UUID记住！！！
func SignUp(p *models.RegisterForm2) (error error) {
	//username := strings.TrimSpace(p.UserName)
	email := strings.TrimSpace(p.Email)
	nickname := strings.TrimSpace(p.NickName)
	// 验证昵称
	if len(nickname) == 0 {
		return errors.New("昵称不能为空")
	}

	// 验证密码
	err := valiadate.IsPassword(p.Password, p.ConfirmPassword)
	if err != nil {
		return err
	}

	// 验证邮箱
	if len(email) > 0 {
		if err := valiadate.IsEmail(email); err != nil {
			return err
		}
		count, err := mysql.CheckEmailExist(p.Email)
		if err != nil {
			if count > 0 {
				return errors.New("邮箱：" + email + " 已被占用")
			} else {
				return err
			}
		}
	} else {
		return errors.New("请输入邮箱")
	}

	//生成UID
	userId, err := snowflake.GetID()
	if err != nil {
		return mysql.ErrorGenIDFailed
	}

	// 构造一个User实例
	u := models.User2{
		UserID:   userId,
		Email:    p.UserName,
		Password: gofunc.EncodePassword(p.Password),
		NickName: p.NickName,
	}
	// 保存进数据库
	return mysql.InsertUser2(u)
}

/**
 * @Author mengjie.han
 * @Description //TODO 判断能否用邮箱登录的逻辑
 * @Date 21:52 2022/2/10
 **/
func Login(p *models.LoginForm2) (user *models.User2, error error) {
	user = &models.User2{
		Email:    p.Email,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT
	//return jwt.GenToken(user.UserID,user.UserName)
	atoken, rtoken, err := jwt.GenToken(user.UserID, user.Email)
	if err != nil {
		return
	}
	user.AccessToken = atoken
	user.RefreshToken = rtoken
	return
}

/**
 * @Author mengjie.han
 * @Description //TODO 修改密码,确定一下输入的邮箱还是什么,后续检查一下内部的sql语句有没有什么错误
 **/
// UpdatePassword 修改密码
func UpdatePassword(p *models.UpdatePasswordForm) (err error) {
	if err := valiadate.IsPassword(p.Password, p.RePassword); err != nil {
		return err
	}
	//通过email在数据库中查询该用户行
	user, err := mysql.GetUserByEmail(p.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("用户不存在")
		}
		return
	}
	if len(user.Password) == 0 {
		return errors.New("你没设置密码，请先设置密码")
	}

	if !gofunc.ValidatePassword(user.Password, p.OldPassword) {
		return errors.New("旧密码验证失败")
	}

	//更新用户行
	return mysql.UpdatePassword(user.Password, user.Email)
}
