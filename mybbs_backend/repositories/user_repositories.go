package repositories

import (
	"gorm.io/gorm"
	"mybbs/model"
)

var UserRepository = newUserRepository()

func newUserRepository() *userRepository {
	return &userRepository{}
}

type userRepository struct {
}

func (r *userRepository) Update(db *gorm.DB, t *model.User) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userRepository) Create(db *gorm.DB, t *model.User) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userRepository) Get(db *gorm.DB, id int64) *model.User {
	ret := &model.User{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) Take(db *gorm.DB, where ...interface{}) *model.User {
	ret := &model.User{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.User{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userRepository) GetByEmail(db *gorm.DB, email string) *model.User {
	return r.Take(db, "email = ?", email)
}

// GetByUsername 根据用户名查找
func (r *userRepository) GetByUsername(db *gorm.DB, username string) *model.User {
	return r.Take(db, "username = ?", username)
}
