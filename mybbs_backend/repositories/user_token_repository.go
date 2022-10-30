package repositories

import (
	"gorm.io/gorm"
	"mybbs/model"
)

var UserTokenRepository = newUserTokenRepository()

func newUserTokenRepository() *userTokenRepository {
	return &userTokenRepository{}
}

type userTokenRepository struct {
}

func (r *userTokenRepository) Create(db *gorm.DB, t *model.UserToken) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userTokenRepository) GetByToken(db *gorm.DB, token string) *model.UserToken {
	if len(token) == 0 {
		return nil
	}
	return r.Take(db, "token = ?", token)
}

func (r *userTokenRepository) Take(db *gorm.DB, where ...interface{}) *model.UserToken {
	ret := &model.UserToken{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}
