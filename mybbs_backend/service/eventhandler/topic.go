package eventhandler

import (
	"gorm.io/gorm"
	"mybbs/model"
	"mybbs/repositories"
)

var TopicService = newTopicService()

func newTopicService() *topicService {
	return &topicService{}
}

type topicService struct{}

// 当帖子被评论的时候，更新最后回复时间、回复数量+1
func (s *topicService) onComment(tx *gorm.DB, topicId int64, comment *model.Comment) error {
	if err := repositories.TopicRepository.Updates(tx, topicId, map[string]interface{}{
		"last_comment_time":    comment.CreateTime,
		"last_comment_user_id": comment.UserId,
		"comment_count":        gorm.Expr("comment_count + 1"),
	}); err != nil {
		return err
	}
	if err := tx.Exec("update t_topic_tag set last_comment_time = ?, last_comment_user_id = ? where topic_id = ?",
		comment.CreateTime, comment.UserId, topicId).Error; err != nil {
		return err
	}
	return nil
}
