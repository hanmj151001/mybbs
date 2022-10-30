package repositories

var TopicRepository = newTopicRepository()

func newTopicRepository() *topicRepository {
	return &topicRepository{}
}

type topicRepository struct {
}

//func (r *topicRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
//	err = db.Model(&model.Topic{}).Where("id = ?", id).Updates(columns).Error
//	return
//}
