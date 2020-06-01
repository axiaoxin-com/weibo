// https://open.weibo.com/wiki/2/statuses/share
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// status	true	string	用户分享到微博的文本内容，必须做URLencode，内容不超过140个汉字，文本中不能包含“#话题词#”，同时文本中必须包含至少一个第三方分享到微博的网页URL，且该URL只能是该第三方（调用方）绑定域下的URL链接，绑定域在“我的应用 － 应用信息 － 基本应用信息编辑 － 安全域名”里设置。
// pic	false	binary	用户想要分享到微博的图片，仅支持JPEG、GIF、PNG图片，上传图片大小限制为<5M。上传图片时，POST方式提交请求，需要采用multipart/form-data编码方式。
// rip	false	string	开发者上报的操作用户真实IP，形如：211.156.0.1。

package weibo

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

// RespStatusesShare 微博成功发送后的返回结构
type RespStatusesShare struct {
	RespError
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

/*StatusesShare 第三方分享一条链接到微博

token 为获取到的access_token内容

status 为微博文字内容

pic 为附带的一张图片，传nil则只发文字
*/
func (w *Weibo) StatusesShare(token, status string, pic io.Reader) (*RespStatusesShare, error) {
	apiURL := "https://api.weibo.com/2/statuses/share.json"
	ip := RealIP()
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	if pic == nil {
		data := url.Values{
			"access_token": {token},
			"status":       {status},
			"rip":          {ip},
		}
		bodyBuf = bytes.NewBufferString(data.Encode())
	} else {
		// close pic if it's a file
		if f, ok := pic.(*os.File); ok {
			defer f.Close()
		}
		picWriter, err := writer.CreateFormFile("pic", "picname.png")
		if err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare CreateFormFile error")
		}
		if _, err := io.Copy(picWriter, pic); err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare io.Copy error")
		}

		if err := writer.WriteField("access_token", token); err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare WriteField access_token error")
		}
		if err := writer.WriteField("status", status); err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare WriteField status error")
		}
		if err := writer.WriteField("rip", ip); err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare WriteField rip error")
		}
		writer.Close() // must close before new request
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, bodyBuf)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShare NewRequest error")
	}
	if pic == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShare Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShare ReadAll error")
	}
	sr := &RespStatusesShare{}
	if err := json.Unmarshal(body, sr); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShare Unmarshal error:"+string(body))
	}
	if sr.Error != "" && sr.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesShare resp error:" + sr.Error)
	}
	if sr.Idstr == "" {
		return nil, errors.New(string(body))
	}
	return sr, nil
}
