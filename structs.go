// 各种结构体定义

package weibo

import (
	"io"
	"net/http"
)

// StatusesShareResp 微博成功发送后的返回结构
type StatusesShareResp struct {
	Visible struct {
		Type   int `json:"type"`
		ListID int `json:"list_id"`
	} `json:"visible"`
	CreatedAt                string        `json:"created_at"`
	ID                       int64         `json:"id"`
	Idstr                    string        `json:"idstr"`
	Mid                      string        `json:"mid"`
	CanEdit                  bool          `json:"can_edit"`
	ShowAdditionalIndication int           `json:"show_additional_indication"`
	Text                     string        `json:"text"`
	TextLength               int           `json:"textLength"`
	SourceAllowclick         int           `json:"source_allowclick"`
	SourceType               int           `json:"source_type"`
	Source                   string        `json:"source"`
	Favorited                bool          `json:"favorited"`
	Truncated                bool          `json:"truncated"`
	InReplyToStatusID        string        `json:"in_reply_to_status_id"`
	InReplyToUserID          string        `json:"in_reply_to_user_id"`
	InReplyToScreenName      string        `json:"in_reply_to_screen_name"`
	PicUrls                  []interface{} `json:"pic_urls"`
	Geo                      interface{}   `json:"geo"`
	IsPaid                   bool          `json:"is_paid"`
	MblogVipType             int           `json:"mblog_vip_type"`
	User                     struct {
		ID               int64  `json:"id"`
		Idstr            string `json:"idstr"`
		Class            int    `json:"class"`
		ScreenName       string `json:"screen_name"`
		Name             string `json:"name"`
		Province         string `json:"province"`
		City             string `json:"city"`
		Location         string `json:"location"`
		Description      string `json:"description"`
		URL              string `json:"url"`
		ProfileImageURL  string `json:"profile_image_url"`
		ProfileURL       string `json:"profile_url"`
		Domain           string `json:"domain"`
		Weihao           string `json:"weihao"`
		Gender           string `json:"gender"`
		FollowersCount   int    `json:"followers_count"`
		FriendsCount     int    `json:"friends_count"`
		PagefriendsCount int    `json:"pagefriends_count"`
		StatusesCount    int    `json:"statuses_count"`
		VideoStatusCount int    `json:"video_status_count"`
		FavouritesCount  int    `json:"favourites_count"`
		CreatedAt        string `json:"created_at"`
		Following        bool   `json:"following"`
		AllowAllActMsg   bool   `json:"allow_all_act_msg"`
		GeoEnabled       bool   `json:"geo_enabled"`
		Verified         bool   `json:"verified"`
		VerifiedType     int    `json:"verified_type"`
		Remark           string `json:"remark"`
		Insecurity       struct {
			SexualContent bool `json:"sexual_content"`
		} `json:"insecurity"`
		Ptype             int    `json:"ptype"`
		AllowAllComment   bool   `json:"allow_all_comment"`
		AvatarLarge       string `json:"avatar_large"`
		AvatarHd          string `json:"avatar_hd"`
		VerifiedReason    string `json:"verified_reason"`
		VerifiedTrade     string `json:"verified_trade"`
		VerifiedReasonURL string `json:"verified_reason_url"`
		VerifiedSource    string `json:"verified_source"`
		VerifiedSourceURL string `json:"verified_source_url"`
		FollowMe          bool   `json:"follow_me"`
		Like              bool   `json:"like"`
		LikeMe            bool   `json:"like_me"`
		OnlineStatus      int    `json:"online_status"`
		BiFollowersCount  int    `json:"bi_followers_count"`
		Lang              string `json:"lang"`
		Star              int    `json:"star"`
		Mbtype            int    `json:"mbtype"`
		Mbrank            int    `json:"mbrank"`
		BlockWord         int    `json:"block_word"`
		BlockApp          int    `json:"block_app"`
		CreditScore       int    `json:"credit_score"`
		UserAbility       int    `json:"user_ability"`
		Urank             int    `json:"urank"`
		StoryReadState    int    `json:"story_read_state"`
		VclubMember       int    `json:"vclub_member"`
		IsTeenager        int    `json:"is_teenager"`
		IsGuardian        int    `json:"is_guardian"`
		IsTeenagerList    int    `json:"is_teenager_list"`
		SpecialFollow     bool   `json:"special_follow"`
		TabManage         string `json:"tab_manage"`
	} `json:"user"`
	RepostsCount         int           `json:"reposts_count"`
	CommentsCount        int           `json:"comments_count"`
	AttitudesCount       int           `json:"attitudes_count"`
	PendingApprovalCount int           `json:"pending_approval_count"`
	IsLongText           bool          `json:"isLongText"`
	RewardExhibitionType int           `json:"reward_exhibition_type"`
	HideFlag             int           `json:"hide_flag"`
	Mlevel               int           `json:"mlevel"`
	BizFeature           int           `json:"biz_feature"`
	HasActionTypeCard    int           `json:"hasActionTypeCard"`
	DarwinTags           []interface{} `json:"darwin_tags"`
	HotWeiboTags         []interface{} `json:"hot_weibo_tags"`
	TextTagTips          []interface{} `json:"text_tag_tips"`
	Mblogtype            int           `json:"mblogtype"`
	UserType             int           `json:"userType"`
	MoreInfoType         int           `json:"more_info_type"`
	PositiveRecomFlag    int           `json:"positive_recom_flag"`
	ContentAuth          int           `json:"content_auth"`
	GifIds               interface{}   `json:"gif_ids"`
	IsShowBulletin       int           `json:"is_show_bulletin"`
	CommentManageInfo    struct {
		CommentManageButton   int `json:"comment_manage_button"`
		CommentPermissionType int `json:"comment_permission_type"`
		ApprovalCommentType   int `json:"approval_comment_type"`
	} `json:"comment_manage_info"`
	PicNum int `json:"pic_num"`
}

// CrackPinFunc 验证码破解方法类型声明
// 验证码图片以 io.Reader 类型传入，返回破解结果字符串
type CrackPinFunc func(io.Reader) (string, error)

// Weibo 实例，在其上实现各类接口
type Weibo struct {
	client        *http.Client
	appkey        string
	appsecret     string
	redirecturi   string
	username      string
	passwd        string
	userAgent     string
	crackPinFuncs []CrackPinFunc
}

// MobileLoginResp 移动端登录的返回结果
type MobileLoginResp struct {
	Retcode int                    `json:"retcode"`
	Msg     string                 `json:"msg"`
	Data    map[string]interface{} `json:"data"`
}

// preLoginResp PC 端 prelogin 的返回结果
type preLoginResp struct {
	Retcode    int    `json:"retcode"`
	Servertime int    `json:"servertime"`
	Pcid       string `json:"pcid"`
	Nonce      string `json:"nonce"`
	Pubkey     string `json:"pubkey"`
	Rsakv      string `json:"rsakv"`
	IsOpenlock int    `json:"is_openlock"`
	Showpin    int    `json:"showpin"`
	Exectime   int    `json:"exectime"`
}

// ssoLoginResp PC 端 ssologin 登录的返回结果
type ssoLoginResp struct {
	Retcode            string   `json:"retcode"`
	Ticket             string   `json:"ticket"`
	UID                string   `json:"uid"`
	Nick               string   `json:"nick"`
	CrossDomainURLList []string `json:"crossDomainUrlList"`
}

// TokenResp 获取 access token 接口的返回结果
type TokenResp struct {
	AccessToken string `json:"access_token"` // access token
	ExpiresIn   int64  `json:"expires_in"`   // ExpiresIn 秒之后token过期
	UID         string `json:"uid"`
	IsRealName  string `json:"isRealName"`
}

// TokenInfoResp 查询 token 信息接口的返回结果
type TokenInfoResp struct {
	UID      string `json:"uid"`
	Appkey   string `json:"appkey"`
	Scope    string `json:"scope"`
	CreateAt string `json:"create_at"`
	ExpireIn string `json:"expire_in"`
}
