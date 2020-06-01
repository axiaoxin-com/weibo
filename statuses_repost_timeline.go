// https://open.weibo.com/wiki/2/statuses/repost_timeline
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// id	true	int64	需要查询的微博ID
// since_id	false	int64	若指定此参数，则返回ID比since_id大的微博（即比since_id时间晚的微博），默认为0。
// max_id	false	int64	若指定此参数，则返回ID小于或等于max_id的微博，默认为0。
// count	false	int	单页返回的记录条数，最大不超过200，默认为20。
// page	false	int	返回结果的页码，默认为1。
// filter_by_author	false	int	作者筛选类型，0：全部、1：我关注的人、2：陌生人，默认为0。

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// RespStatusesRepostTimeline StatusesRepostTimeline 接口返回结构
type RespStatusesRepostTimeline struct {
	RespError
	Reposts []struct {
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
		RetweetedStatus struct {
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
		} `json:"retweeted_status"`
	} `json:"reposts"`
	PreviousCursor int `json:"previous_cursor"`
	NextCursor     int `json:"next_cursor"`
	TotalNumber    int `json:"total_number"`
}

// StatusesRepostTimeline 获取指定微博的转发微博列表
// id 微博ID
// sinceID 返回ID比sinceID大的微博（即比since_id时间晚的微博）
// maxID 返回ID小于或等于max_id的微博0。
// count 单页返回的记录条数，最大不超过200。
// page	返回结果的页码。
// filterByAuthor	作者筛选类型，0：全部、1：我关注的人、2：陌生人
func (w *Weibo) StatusesRepostTimeline(token string, id, sinceID, maxID int64, count, page, filterByAuthor int) (*RespStatusesRepostTimeline, error) {
	apiURL := "https://api.weibo.com/2/statuses/repost_timeline.json"
	data := url.Values{
		"access_token":     {token},
		"id":               {strconv.FormatInt(id, 10)},
		"since_id":         {strconv.FormatInt(sinceID, 10)},
		"max_id":           {strconv.FormatInt(maxID, 10)},
		"count":            {strconv.Itoa(count)},
		"page":             {strconv.Itoa(page)},
		"filter_by_author": {strconv.Itoa(filterByAuthor)},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesRepostTimeline NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesRepostTimeline Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesRepostTimeline ReadAll error")
	}
	r := &RespStatusesRepostTimeline{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesRepostTimeline Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesRepostTimeline resp error:" + r.Error)
	}
	return r, nil
}
