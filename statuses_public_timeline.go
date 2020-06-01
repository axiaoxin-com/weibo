// https://open.weibo.com/wiki/2/statuses/public_timeline
// 请求参数
//  count	false	int	单页返回的记录条数，默认为50。
//  page	false	int	返回结果的页码，默认为1。
//  base_app	false	int	是否只获取当前应用的数据。0为否（所有数据），1为是（仅当前应用），默认为0。reply_comment	object	评论来源评论，当本评论属于对另一评论的回复时返回此字段

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// RespStatusesPublicTimeline StatusesPublicTimeline 接口返回结果
type RespStatusesPublicTimeline struct {
	RespError
	Statuses []struct {
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
	} `json:"statuses"`
	PreviousCursor int   `json:"previous_cursor"`
	NextCursor     int64 `json:"next_cursor"`
	TotalNumber    int   `json:"total_number"`
}

// StatusesPublicTimeline 返回最新的公共微博
// count 单页返回的记录条数
// page 返回结果的页码
// baseApp 是否只获取当前应用的数据
func (w *Weibo) StatusesPublicTimeline(token string, count, page, baseApp int) (*RespStatusesPublicTimeline, error) {
	apiURL := "https://api.weibo.com/2/statuses/public_timeline.json"
	data := url.Values{
		"access_token": {token},
		"count":        {strconv.Itoa(count)},
		"page":         {strconv.Itoa(page)},
		"base_app":     {strconv.Itoa(baseApp)},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesPublicTimeline NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesPublicTimeline Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesPublicTimeline ReadAll error")
	}
	r := &RespStatusesPublicTimeline{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesPublicTimeline Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesPublicTimeline resp error:" + r.Error)
	}
	return r, nil
}
