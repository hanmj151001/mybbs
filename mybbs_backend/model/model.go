package model

import (
	"database/sql"
	"time"
)

//目录管理方式有问题，需要解决一下

var Models = []interface{}{
	&User{}, &UserToken{}, &Article{}, &Comment{}, &SysConfig{},
}

type Model struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}

type User struct {
	Model
	Username      sql.NullString `gorm:"size:32;unique;" json:"username" form:"username"`                  // 用户名
	Email         sql.NullString `gorm:"size:128;unique;" json:"email" form:"email"`                       // 邮箱
	EmailVerified bool           `gorm:"not null;default:false" json:"emailVerified" form:"emailVerified"` // 邮箱是否验证
	Nickname      string         `gorm:"size:16;" json:"nickname" form:"nickname"`                         // 昵称
	Avatar        string         `gorm:"type:text" json:"avatar" form:"avatar"`                            // 头像
	//Gender           constants.Gender `gorm:"size:16;default:''" json:"gender" form:"gender"`                     // 性别
	Birthday         *time.Time `json:"birthday" form:"birthday"`                                           // 生日
	BackgroundImage  string     `gorm:"type:text" json:"backgroundImage" form:"backgroundImage"`            // 个人中心背景图片
	Password         string     `gorm:"size:512" json:"password" form:"password"`                           // 密码
	HomePage         string     `gorm:"size:1024" json:"homePage" form:"homePage"`                          // 个人主页
	Description      string     `gorm:"type:text" json:"description" form:"description"`                    // 个人描述
	Score            int        `gorm:"not null;index:idx_user_score" json:"score" form:"score"`            // 积分
	Status           int        `gorm:"index:idx_user_status;not null" json:"status" form:"status"`         // 状态
	TopicCount       int        `gorm:"not null" json:"topicCount" form:"topicCount"`                       // 帖子数量
	CommentCount     int        `gorm:"not null" json:"commentCount" form:"commentCount"`                   // 跟帖数量
	FollowCount      int        `gorm:"not null" json:"followCount" form:"followCount"`                     // 关注数量
	FansCount        int        `gorm:"not null" json:"fansCount" form:"fansCount"`                         // 粉丝数量
	Roles            string     `gorm:"type:text" json:"roles" form:"roles"`                                // 角色
	ForbiddenEndTime int64      `gorm:"not null;default:0" json:"forbiddenEndTime" form:"forbiddenEndTime"` // 禁言结束时间
	CreateTime       int64      `json:"createTime" form:"createTime"`                                       // 创建时间
	UpdateTime       int64      `json:"updateTime" form:"updateTime"`                                       // 更新时间
}

type UserToken struct {
	Model
	Token      string `gorm:"size:32;unique;not null" json:"token" form:"token"`
	UserId     int64  `gorm:"not null;index:idx_user_token_user_id;" json:"userId" form:"userId"`
	ExpiredAt  int64  `gorm:"not null" json:"expiredAt" form:"expiredAt"`
	Status     int    `gorm:"not null;index:idx_user_token_status" json:"status" form:"status"`
	CreateTime int64  `gorm:"not null" json:"createTime" form:"createTime"`
}

// 文章
type Article struct {
	Model
	UserId      int64  `gorm:"index:idx_article_user_id" json:"userId" form:"userId"`             // 所属用户编号
	Title       string `gorm:"size:128;not null;" json:"title" form:"title"`                      // 标题
	Summary     string `gorm:"type:text" json:"summary" form:"summary"`                           // 摘要
	Content     string `gorm:"type:longtext;not null;" json:"content" form:"content"`             // 内容
	ContentType string `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`   // 内容类型：markdown、html
	Status      int    `gorm:"type:int(11);index:idx_article_status" json:"status" form:"status"` // 状态
	SourceUrl   string `gorm:"type:text" json:"sourceUrl" form:"sourceUrl"`                       // 原文链接
	ViewCount   int64  `gorm:"not null;index:idx_view_count;" json:"viewCount" form:"viewCount"`  // 查看数量
	CreateTime  int64  `gorm:"index:idx_article_create_time" json:"createTime" form:"createTime"` // 创建时间
	UpdateTime  int64  `json:"updateTime" form:"updateTime"`                                      // 更新时间
}

// 评论
type Comment struct {
	Model
	UserId       int64  `gorm:"index:idx_comment_user_id;not null" json:"userId" form:"userId"`             // 用户编号
	EntityType   string `gorm:"index:idx_comment_entity_type;not null" json:"entityType" form:"entityType"` // 被评论实体类型
	EntityId     int64  `gorm:"index:idx_comment_entity_id;not null" json:"entityId" form:"entityId"`       // 被评论实体编号
	Content      string `gorm:"type:text;not null" json:"content" form:"content"`                           // 内容
	ImageList    string `gorm:"type:longtext" json:"imageList" form:"imageList"`                            // 图片
	ContentType  string `gorm:"type:varchar(32);not null" json:"contentType" form:"contentType"`            // 内容类型：markdown、html
	QuoteId      int64  `gorm:"not null"  json:"quoteId" form:"quoteId"`                                    // 引用的评论编号
	LikeCount    int64  `gorm:"not null;default:0" json:"likeCount" form:"likeCount"`                       // 点赞数量
	CommentCount int64  `gorm:"not null;default:0" json:"commentCount" form:"commentCount"`                 // 评论数量
	UserAgent    string `gorm:"size:1024" json:"userAgent" form:"userAgent"`                                // UserAgent
	Ip           string `gorm:"size:128" json:"ip" form:"ip"`                                               // IP
	Status       int    `gorm:"int;index:idx_comment_status" json:"status" form:"status"`                   // 状态：0：待审核、1：审核通过、2：审核失败、3：已发布
	CreateTime   int64  `json:"createTime" form:"createTime"`                                               // 创建时间
}

// 系统配置
type SysConfig struct {
	Model
	Key         string `gorm:"not null;size:128;unique" json:"key" form:"key"` // 配置key
	Value       string `gorm:"type:text" json:"value" form:"value"`            // 配置值
	Name        string `gorm:"not null;size:32" json:"name" form:"name"`       // 配置名称
	Description string `gorm:"size:128" json:"description" form:"description"` // 配置描述
	CreateTime  int64  `gorm:"not null" json:"createTime" form:"createTime"`   // 创建时间
	UpdateTime  int64  `gorm:"not null" json:"updateTime" form:"updateTime"`   // 更新时间
}
