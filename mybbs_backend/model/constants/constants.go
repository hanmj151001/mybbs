package constants

const (
	DefaultTokenExpireDays       = 3   // 用户登录token默认有效期
	SummaryLen                   = 256 // 摘要长度
	UploadMaxM                   = 10
	UploadMaxBytes         int64 = 1024 * 1024 * 1024 * UploadMaxM
)

// 系统配置
const (
	SysConfigSiteTitle        = "siteTitle"        // 站点标题
	SysConfigSiteDescription  = "siteDescription"  // 站点描述
	SysConfigSiteKeywords     = "siteKeywords"     // 站点关键字
	SysConfigSiteNavs         = "siteNavs"         // 站点导航
	SysConfigSiteNotification = "siteNotification" // 站点公告
	SysConfigRecommendTags    = "recommendTags"    // 推荐标签
	SysConfigUrlRedirect      = "urlRedirect"      // 是否开启链接跳转
	SysConfigScoreConfig      = "scoreConfig"      // 分数配置
	SysConfigDefaultNodeId    = "defaultNodeId"    // 发帖默认节点
	SysConfigArticlePending   = "articlePending"   // 是否开启文章审核

	SysConfigTokenExpireDays   = "tokenExpireDays"   // 登录Token有效天数
	SysConfigEnableHideContent = "enableHideContent" // 启用回复可见功能
)

// EntityType
const (
	EntityArticle = "article"
	EntityTopic   = "topic"
	EntityComment = "comment"
	EntityUser    = "user"
	EntityCheckIn = "checkIn"
)

// 用户角色
const (
	RoleOwner = "owner" // 站长
	RoleAdmin = "admin" // 管理员
	RoleUser  = "user"  // 用户
)

// 操作类型
const (
	OpTypeCreate          = "create"
	OpTypeDelete          = "delete"
	OpTypeUpdate          = "update"
	OpTypeForbidden       = "forbidden"
	OpTypeRemoveForbidden = "removeForbidden"
)

// 状态
const (
	StatusOk      = 0 // 正常
	StatusDeleted = 1 // 删除
	StatusPending = 2 // 待审核
)

// 用户类型
const (
	UserTypeNormal = 0 // 普通用户
	UserTypeGzh    = 1 // 公众号用户
)

// 内容类型
const (
	ContentTypeHtml     = "html"
	ContentTypeMarkdown = "markdown"
	ContentTypeText     = "text"
)

// 第三方账号类型
const (
	ThirdAccountTypeGithub = "github"
	ThirdAccountTypeOSC    = "osc"
	ThirdAccountTypeQQ     = "qq"
)

// 积分操作类型
const (
	ScoreTypeIncr = 0 // 积分+
	ScoreTypeDecr = 1 // 积分-
)

type TopicType int

const (
	TopicTypeTopic TopicType = 0
	TopicTypeTweet TopicType = 1
)

type LoginMethod string

const (
	LoginMethodQQ       LoginMethod = "qq"
	LoginMethodGithub   LoginMethod = "github"
	LoginMethodPassword LoginMethod = "password"
)

const (
	FollowStatusNONE   = 0
	FollowStatusFollow = 1
	FollowStatusBoth   = 2
)

const (
	NodeIdNewest    int64 = 0
	NodeIdRecommend int64 = -1
	NodeIdFollow    int64 = -2
)

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)
