package repositories

import (
	"gorm.io/gorm"
	"mybbs/model"
)

var CommentRepository = newCommentRepository()

func newCommentRepository() *commentRepository {
	return &commentRepository{}
}

type commentRepository struct {
}

func (r *commentRepository) Create(db *gorm.DB, t *model.Comment) (err error) {
	err = db.Create(t).Error
	return
}
