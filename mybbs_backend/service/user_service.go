package service

import (
	"errors"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"mybbs/model"
	"mybbs/model/constants"
	"mybbs/pkg/gofunc"
	"mybbs/pkg/validate"
	"mybbs/repositories"
	"strings"
)

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (s *userService) Get(id int64) *model.User {
	return repositories.UserRepository.Get(sqls.DB(), id)
}

func (s *userService) Take(where ...interface{}) *model.User {
	return repositories.UserRepository.Take(sqls.DB(), where...)
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.UserRepository.UpdateColumn(sqls.DB(), id, name, value)
	//cache.UserCache.Invalidate(id)
	return err
}

// GetByEmail 根据邮箱查找
func (s *userService) GetByEmail(email string) *model.User {
	return repositories.UserRepository.GetByEmail(sqls.DB(), email)
}

// GetByUsername 根据用户名查找
func (s *userService) GetByUsername(username string) *model.User {
	return repositories.UserRepository.GetByUsername(sqls.DB(), username)
}

// isUsernameExists 用户名是否存在
func (s *userService) isUsernameExists(username string) bool {
	return s.GetByUsername(username) != nil
}

// SignUp 注册
func (s *userService) SignUp(username, email, nickname, password, rePassword string) (*model.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	nickname = strings.TrimSpace(nickname)

	// 验证昵称
	if len(nickname) == 0 {
		return nil, errors.New("昵称不能为空")
	}

	// 验证密码
	err := validate.IsPassword(password, rePassword)
	if err != nil {
		return nil, err
	}

	// 验证邮箱
	if len(email) > 0 {
		if err := validate.IsEmail(email); err != nil {
			return nil, err
		}
		if s.GetByEmail(email) != nil {
			return nil, errors.New("邮箱：" + email + " 已被占用")
		}
	} else {
		return nil, errors.New("请输入邮箱")
	}

	// 验证用户名
	if len(username) > 0 {
		//if err := validate.IsUsername(username); err != nil {
		//	return nil, err
		//}
		if s.isUsernameExists(username) {
			return nil, errors.New("用户名：" + username + " 已被占用")
		}
	}

	user := &model.User{
		Username:   sqls.SqlNullString(username),
		Email:      sqls.SqlNullString(email),
		Nickname:   nickname,
		Password:   gofunc.EncodePassword(password),
		Status:     constants.StatusOk,
		CreateTime: dates.NowTimestamp(),
		UpdateTime: dates.NowTimestamp(),
	}

	err = repositories.UserRepository.Create(sqls.DB(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SignIn 登录
func (s *userService) SignIn(username, password string) (*model.User, error) {
	if len(username) == 0 {
		return nil, errors.New("用户名/邮箱不能为空")
	}
	if len(password) == 0 {
		return nil, errors.New("密码不能为空")
	}
	var user *model.User = nil
	if err := validate.IsEmail(username); err == nil { // 如果用户输入的是邮箱
		user = s.GetByEmail(username)
	} else {
		user = s.GetByUsername(username)
	}
	if user == nil || user.Status != constants.StatusOk {
		return nil, errors.New("用户不存在或被禁用")
	}
	if !gofunc.ValidatePassword(user.Password, password) {
		return nil, errors.New("密码错误")
	}
	return user, nil
}

// UpdatePassword 修改密码
func (s *userService) UpdatePassword(userId int64, oldPassword, password, rePassword string) error {
	if err := validate.IsPassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)

	if len(user.Password) == 0 {
		return errors.New("你没设置密码，请先设置密码")
	}

	if !gofunc.ValidatePassword(user.Password, oldPassword) {
		return errors.New("旧密码验证失败")
	}

	return s.UpdateColumn(userId, "password", gofunc.EncodePassword(password))
}
