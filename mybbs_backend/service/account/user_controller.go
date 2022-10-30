package account

import (
	"errors"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/sqls"
	//"github.com/mlogclub/simple/web"
	"mybbs/model"
	"mybbs/model/constants"
	"mybbs/pkg/validate"
	"mybbs/repositories"

	"strings"
)

// 邮箱验证邮件有效期（小时）
const emailVerifyExpireHour = 24

var UserService = newUserService()

func newUserService() *userService {
	return &userService{}
}

type userService struct {
}

func (s *userService) Get(id int64) *model.User {
	return repositories.UserRepository.Get(sqls.DB(), id)
}

func (s *userService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.UserRepository.UpdateColumn(sqls.DB(), id, name, value)
	cache.UserCache.Invalidate(id)
	return err
}

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
		Password:   passwd.EncodePassword(password),
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
	if !passwd.ValidatePassword(user.Password, password) {
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

	if !passwd.ValidatePassword(user.Password, oldPassword) {
		return errors.New("旧密码验证失败")
	}

	return s.UpdateColumn(userId, "password", passwd.EncodePassword(password))
}

// CheckPostStatus 用于在发表内容时检查用户状态
//func (s *userService) CheckPostStatus(user *model.User) error {
//	if user == nil {
//		return errs.NotLogin
//	}
//	if user.Status != constants.StatusOk {
//		return errs.UserDisabled
//	}
//	if user.IsForbidden() {
//		return errs.ForbiddenError
//	}
//	observeSeconds := SysConfigService.GetInt(constants.SysConfigUserObserveSeconds)
//	if user.InObservationPeriod(observeSeconds) {
//		return web.NewError(errs.InObservationPeriod.Code, "账号尚在观察期，观察期时长："+strconv.Itoa(observeSeconds)+"秒，请稍后再试")
//	}
//	return nil
//}

// isEmailExists 邮箱是否存在
func (s *userService) isEmailExists(email string) bool {
	if len(email) == 0 { // 如果邮箱为空，那么就认为是不存在
		return false
	}
	return s.GetByEmail(email) != nil
}

// isUsernameExists 用户名是否存在
func (s *userService) isUsernameExists(username string) bool {
	return s.GetByUsername(username) != nil
}

// UpdateAvatar 更新头像
func (s *userService) UpdateAvatar(userId int64, avatar string) error {
	return s.UpdateColumn(userId, "avatar", avatar)
}

// UpdateNickname 更新昵称
func (s *userService) UpdateNickname(userId int64, nickname string) error {
	return s.UpdateColumn(userId, "nickname", nickname)
}

// SetPassword 设置密码
func (s *userService) SetPassword(userId int64, password, rePassword string) error {
	if err := validate.IsPassword(password, rePassword); err != nil {
		return err
	}
	user := s.Get(userId)
	if len(user.Password) > 0 {
		return errors.New("你已设置了密码，如需修改请前往修改页面。")
	}
	password = passwd.EncodePassword(password)
	return s.UpdateColumn(userId, "password", password)
}

// GetByEmail 根据邮箱查找
func (s *userService) GetByEmail(email string) *model.User {
	return repositories.UserRepository.GetByEmail(sqls.DB(), email)
}

// GetByUsername 根据用户名查找
func (s *userService) GetByUsername(username string) *model.User {
	return repositories.UserRepository.GetByUsername(sqls.DB(), username)
}
