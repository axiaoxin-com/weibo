// https://open.weibo.com/wiki/2/comments/show_batch
// 请求参数
//   access_token	true	string	采用 OAuth 授权方式为必填参数， OAuth 授权后获得。
//   cids	true	int64	需要查询的批量评论 ID ，用半角逗号分隔，最大 50

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// RespCommentsShowBatch CommentsShowBatch 接口返回结构
type RespCommentsShowBatch []struct {
	CreatedAt string `json:"created_at"`
	ID        int64  `json:"id"`
	Text      string `json:"text"`
	Source    string `json:"source"`
	Mid       string `json:"mid"`
	User      struct {
		ID               int    `json:"id"`
		ScreenName       string `json:"screen_name"`
		Name             string `json:"name"`
		Province         string `json:"province"`
		City             string `json:"city"`
		Location         string `json:"location"`
		Description      string `json:"description"`
		URL              string `json:"url"`
		ProfileImageURL  string `json:"profile_image_url"`
		Domain           string `json:"domain"`
		Gender           string `json:"gender"`
		FollowersCount   int    `json:"followers_count"`
		FriendsCount     int    `json:"friends_count"`
		StatusesCount    int    `json:"statuses_count"`
		FavouritesCount  int    `json:"favourites_count"`
		CreatedAt        string `json:"created_at"`
		Following        bool   `json:"following"`
		AllowAllActMsg   bool   `json:"allow_all_act_msg"`
		Remark           string `json:"remark"`
		GeoEnabled       bool   `json:"geo_enabled"`
		Verified         bool   `json:"verified"`
		AllowAllComment  bool   `json:"allow_all_comment"`
		AvatarLarge      string `json:"avatar_large"`
		VerifiedReason   string `json:"verified_reason"`
		FollowMe         bool   `json:"follow_me"`
		OnlineStatus     int    `json:"online_status"`
		BiFollowersCount int    `json:"bi_followers_count"`
	} `json:"user"`
	Status struct {
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
		RepostsCount        int           `json:"reposts_count"`
		CommentsCount       int           `json:"comments_count"`
		Annotations         []interface{} `json:"annotations"`
		User                struct {
			ID               int    `json:"id"`
			ScreenName       string `json:"screen_name"`
			Name             string `json:"name"`
			Province         string `json:"province"`
			City             string `json:"city"`
			Location         string `json:"location"`
			Description      string `json:"description"`
			URL              string `json:"url"`
			ProfileImageURL  string `json:"profile_image_url"`
			Domain           string `json:"domain"`
			Gender           string `json:"gender"`
			FollowersCount   int    `json:"followers_count"`
			FriendsCount     int    `json:"friends_count"`
			StatusesCount    int    `json:"statuses_count"`
			FavouritesCount  int    `json:"favourites_count"`
			CreatedAt        string `json:"created_at"`
			Following        bool   `json:"following"`
			AllowAllActMsg   bool   `json:"allow_all_act_msg"`
			Remark           string `json:"remark"`
			GeoEnabled       bool   `json:"geo_enabled"`
			Verified         bool   `json:"verified"`
			AllowAllComment  bool   `json:"allow_all_comment"`
			AvatarLarge      string `json:"avatar_large"`
			VerifiedReason   string `json:"verified_reason"`
			FollowMe         bool   `json:"follow_me"`
			OnlineStatus     int    `json:"online_status"`
			BiFollowersCount int    `json:"bi_followers_count"`
		} `json:"user"`
	} `json:"status"`
}

// CommentsShowBatch 根据评论 ID 批量返回评论信息
// cids 需要查询的批量评论 ID
func (w *Weibo) CommentsShowBatch(token string, cids ...int64) (*RespCommentsShowBatch, error) {
	apiURL := "https://api.weibo.com/2/comments/show_batch.json"
	sCids := []string{}
	for _, cid := range cids {
		sCids = append(sCids, strconv.FormatInt(cid, 10))
	}
	cidsStr := strings.Join(sCids, ",")
	data := url.Values{
		"access_token": {token},
		"cids":         {cidsStr},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsShowBatch NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsShowBatch Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsShowBatch ReadAll error")
	}
	e := &RespError{}
	json.Unmarshal(body, e)
	if e.Error != "" && e.ErrorCode != 0 {
		return nil, errors.New("weibo CommentsShowBatch resp error:" + e.Error)
	}
	r := &RespCommentsShowBatch{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo CommentsShowBatch Unmarshal error:"+string(body))
	}
	return r, nil
}
