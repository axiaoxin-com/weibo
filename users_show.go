// https://open.weibo.com/wiki/2/users/show
// 请求参数
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// uid	false	int64	需要查询的用户ID。
// screen_name	false	string	需要查询的用户昵称。
//
// 注意事项
// 参数uid与screen_name二者必选其一，且只能选其一；
// 接口升级后，对未授权本应用的uid，将无法获取其个人简介、认证原因、粉丝数、关注数、微博数及最近一条微博内容。
//
// 返回字段说明
// id	int64	用户UID
// idstr	string	字符串型的用户UID
// screen_name	string	用户昵称
// name	string	友好显示名称
// province	int	用户所在省级ID
// city	int	用户所在城市ID
// location	string	用户所在地
// description	string	用户个人描述
// url	string	用户博客地址
// profile_image_url	string	用户头像地址（中图），50×50像素
// profile_url	string	用户的微博统一URL地址
// domain	string	用户的个性化域名
// weihao	string	用户的微号
// gender	string	性别，m：男、f：女、n：未知
// followers_count	int	粉丝数
// friends_count	int	关注数
// statuses_count	int	微博数
// favourites_count	int	收藏数
// created_at	string	用户创建（注册）时间
// following	boolean	暂未支持
// allow_all_act_msg	boolean	是否允许所有人给我发私信，true：是，false：否
// geo_enabled	boolean	是否允许标识用户的地理位置，true：是，false：否
// verified	boolean	是否是微博认证用户，即加V用户，true：是，false：否
// verified_type	int	暂未支持
// remark	string	用户备注信息，只有在查询用户关系时才返回此字段
// status	object	用户的最近一条微博信息字段 详细
// allow_all_comment	boolean	是否允许所有人对我的微博进行评论，true：是，false：否
// avatar_large	string	用户头像地址（大图），180×180像素
// avatar_hd	string	用户头像地址（高清），高清头像原图
// verified_reason	string	认证原因
// follow_me	boolean	该用户是否关注当前登录用户，true：是，false：否
// online_status	int	用户的在线状态，0：不在线、1：在线
// bi_followers_count	int	用户的互粉数
// lang	string	用户当前的语言版本，zh-cn：简体中文，zh-tw：繁体中文，en：英语

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// RespUsersShow UsersShow 接口返回结构
type RespUsersShow struct {
	RespError
	ID              int    `json:"id"`
	ScreenName      string `json:"screen_name"`
	Name            string `json:"name"`
	Province        string `json:"province"`
	City            string `json:"city"`
	Location        string `json:"location"`
	Description     string `json:"description"`
	URL             string `json:"url"`
	ProfileImageURL string `json:"profile_image_url"`
	Domain          string `json:"domain"`
	Gender          string `json:"gender"`
	FollowersCount  int    `json:"followers_count"`
	FriendsCount    int    `json:"friends_count"`
	StatusesCount   int    `json:"statuses_count"`
	FavouritesCount int    `json:"favourites_count"`
	CreatedAt       string `json:"created_at"`
	Following       bool   `json:"following"`
	AllowAllActMsg  bool   `json:"allow_all_act_msg"`
	GeoEnabled      bool   `json:"geo_enabled"`
	Verified        bool   `json:"verified"`
	Status          struct {
		CreatedAt           string        `json:"created_at"`
		ID                  int64         `json:"id"`
		Text                string        `json:"text"`
		Source              string        `json:"source"`
		Favorited           bool          `json:"favorited"`
		Truncated           bool          `json:"truncated"`
		InReplyToStatusID   string        `json:"in_reply_to_status_id"`
		InReplyToUserID     string        `json:"in_reply_to_user_id"`
		InReplyToScreenName string        `json:"in_reply_to_screen_name"`
		Geo                 interface{}   `json:"geo"`
		Mid                 string        `json:"mid"`
		Annotations         []interface{} `json:"annotations"`
		RepostsCount        int           `json:"reposts_count"`
		CommentsCount       int           `json:"comments_count"`
	} `json:"status"`
	AllowAllComment  bool   `json:"allow_all_comment"`
	AvatarLarge      string `json:"avatar_large"`
	VerifiedReason   string `json:"verified_reason"`
	FollowMe         bool   `json:"follow_me"`
	OnlineStatus     int    `json:"online_status"`
	BiFollowersCount int    `json:"bi_followers_count"`
}

// UsersShow 根据用户ID获取用户信息
func (w *Weibo) UsersShow(token string, uid int64, screenName string) (*RespUsersShow, error) {
	apiURL := "https://api.weibo.com/2/users/show.json"
	data := url.Values{
		"access_token": {token},
	}
	if uid > 0 {
		data.Add("uid", strconv.FormatInt(uid, 10))
	}
	if screenName != "" {
		data.Add("screen_name", screenName)
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo UsersShow NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo UsersShow Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo UsersShow ReadAll error")
	}
	r := &RespUsersShow{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo UsersShow Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo UsersShow resp error:" + r.Error)
	}
	return r, nil
}
