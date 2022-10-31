package valiadate

import (
	"bluebell_backend/pkg/gofunc"
	"bluebell_backend/pkg/strs"
	"errors"
	"regexp"
)

// IsPassword 是否是合法的密码
func IsPassword(password, rePassword string) error {
	if gofunc.IsBlank(password) {
		return errors.New("请输入密码")
	}
	if strs.RuneLen(password) < 6 {
		return errors.New("密码格式简单")
	}
	if password != rePassword {
		return errors.New("两次输入密码不一致")
	}
	return nil
}

// IsEmail 正则表达式验证是否是合法的邮箱
func IsEmail(email string) (err error) {
	if gofunc.IsBlank(email) {
		err = errors.New("邮箱格式不符合规范")
		return
	}
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		err = errors.New("邮箱格式不符合规范")
	}
	return
}
