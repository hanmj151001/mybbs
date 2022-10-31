package gofunc

import "golang.org/x/crypto/bcrypt"

//EncodePassword 对密码进行加密
func EncodePassword(rawPassword string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return string(hash)
}

func ValidatePassword(encodePassword, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encodePassword), []byte(inputPassword))
	return err == nil
}
