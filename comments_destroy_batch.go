// https://open.weibo.com/wiki/2/comments/destroy_batch
// 请求参数
//   access_token	true	string	采用 OAuth 授权方式为必填参数， OAuth 授权后获得。
//   cids	true	int64	需要删除的评论 ID ，用半角逗号隔开，最多 20 个。

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

// RespCommentsDestroyBatch CommentsDestroyBatch 接口返回结构
type RespCommentsDestroyBatch []struct {
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

// CommentsDestroyBatch 根据评论 ID 批量删除评论
// cids 需要删除的评论 ID
func (w *Weibo) CommentsDestroyBatch(token string, cids ...int64) (*RespCommentsDestroyBatch, error) {
	apiURL := "https://api.weibo.com/2/comments/destroy_batch.json"
	sCids := []string{}
	for _, cid := range cids {
		sCids = append(sCids, strconv.FormatInt(cid, 10))
	}
	cidsStr := strings.Join(sCids, ",")
	data := url.Values{
		"access_token": {token},
		"cids":         {cidsStr},
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsDestroyBatch NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsDestroyBatch Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsDestroyBatch ReadAll error")
	}
	e := &RespError{}
	json.Unmarshal(body, e)
	if e.Error != "" && e.ErrorCode != 0 {
		return nil, errors.New("weibo CommentsDestroyBatch resp error:" + e.Error)
	}
	r := &RespCommentsDestroyBatch{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo CommentsDestroyBatch Unmarshal error:"+string(body))
	}
	return r, nil
}
