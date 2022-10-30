package eventhandler

import (
	"errors"
	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"mybbs/model"
	"mybbs/model/constants"
	"mybbs/pkg/gofunc"
	"strings"
)

var ArticleService = newArticleService()

func newArticleService() *articleService {
	return &articleService{}
}

type articleService struct {
}

// 发布文章
func (s *articleService) Publish(userId int64, form model.CreateArticleForm) (article *model.Article, err error) {
	form.Title = strings.TrimSpace(form.Title)
	form.Summary = strings.TrimSpace(form.Summary)
	form.Content = strings.TrimSpace(form.Content)

	if gofunc.IsBlank(form.Title) {
		return nil, errors.New("标题不能为空")
	}
	if gofunc.IsBlank(form.Content) {
		return nil, errors.New("内容不能为空")
	}

	// 获取后台配置 否是开启发表文章审核
	status := constants.StatusOk
	sysConfigArticlePending := cache.SysConfigCache.GetValue(constants.SysConfigArticlePending)
	if strings.ToLower(sysConfigArticlePending) == "true" {
		status = constants.StatusPending
	}

	article = &model.Article{
		UserId:      userId,
		Title:       form.Title,
		Summary:     form.Summary,
		Content:     form.Content,
		ContentType: form.ContentType,
		Status:      status,
		SourceUrl:   form.SourceUrl,
		CreateTime:  dates.NowTimestamp(),
		UpdateTime:  dates.NowTimestamp(),
	}

	err = sqls.DB().Transaction(func(tx *gorm.DB) error {
		tagIds := repositories.TagRepository.GetOrCreates(tx, form.Tags)
		err := repositories.ArticleRepository.Create(tx, article)
		if err != nil {
			return err
		}
		repositories.ArticleTagRepository.AddArticleTags(tx, article.Id, tagIds)
		return nil
	})

	if err == nil {
		seo.Push(bbsurls.ArticleUrl(article.Id))
	}
	return
}

// 倒序扫描
func (s *articleService) ScanDesc(callback func(articles []model.Article)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ArticleRepository.Find(sqls.DB(), sqls.NewCnd().
			Cols("id", "status", "create_time", "update_time").
			Lt("id", cursor).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}

// 倒序扫描
func (s *articleService) ScanDescWithDate(dateFrom, dateTo int64, callback func(articles []model.Article)) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.ArticleRepository.Find(sqls.DB(), sqls.NewCnd().
			Cols("id", "status", "create_time", "update_time").
			Lt("id", cursor).Gte("create_time", dateFrom).Lt("create_time", dateTo).Desc("id").Limit(1000))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}
