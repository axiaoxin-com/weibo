// https://open.weibo.com/wiki/2/statuses/user_timeline
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// uid	false	int64	需要查询的用户ID。
// screen_name	false	string	需要查询的用户昵称。
// since_id	false	int64	若指定此参数，则返回ID比since_id大的微博（即比since_id时间晚的微博），默认为0。
// max_id	false	int64	若指定此参数，则返回ID小于或等于max_id的微博，默认为0。
// count	false	int	单页返回的记录条数，最大不超过100，超过100以100处理，默认为20。
// page	false	int	返回结果的页码，默认为1。
// base_app	false	int	是否只获取当前应用的数据。0为否（所有数据），1为是（仅当前应用），默认为0。
// feature	false	int	过滤类型ID，0：全部、1：原创、2：图片、3：视频、4：音乐，默认为0。
// trim_user	false	int	返回值中user字段开关，0：返回完整user字段、1：user字段仅返回user_id，默认为0。
//
// 获取自己的微博，参数uid与screen_name可以不填，则自动获取当前登录用户的微博；
// 指定获取他人的微博，参数uid与screen_name二者必选其一，且只能选其一；
// 接口升级后：uid与screen_name只能为当前授权用户；
// 读取当前授权用户所有关注人最新微博列表，请使用：获取当前授权用户及其所关注用户的最新微博接口（statuses/home_timeline）；
// 此接口最多只返回最新的5条数据，官方移动SDK调用可返回10条；

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// StatusesUserTimeline 获取当前授权用户最新发表的微博列表
// uid	int64	需要查询的用户ID。
// screenName	string	需要查询的用户昵称。
// sinceID	int64	返回ID比since_id大的微博（即比since_id时间晚的微博）。
// maxID	int64	返回ID小于或等于max_id的微博。
// count	int	单页返回的记录条数，最大不超过100，超过100以100处理。
// page	int	返回结果的页码
// baseApp	int	是否只获取当前应用的数据。0为否（所有数据），1为是（仅当前应用）。
// feature	int	过滤类型ID，0：全部、1：原创、2：图片、3：视频、4：音乐。
// trimUser	int	返回值中user字段开关，0：返回完整user字段、1：user字段仅返回user_id。
func (w *Weibo) StatusesUserTimeline(token string, uid int64, screenName string, sinceID, maxID int64, count, page, baseApp, feature, trimUser int) (*StatusesUserTimelineResp, error) {
	apiURL := "https://api.weibo.com/2/statuses/user_timeline.json"
	data := url.Values{
		"access_token": {token},
		"uid":          {strconv.FormatInt(uid, 10)},
		"screen_name":  {screenName},
		"since_id":     {strconv.FormatInt(sinceID, 10)},
		"max_id":       {strconv.FormatInt(maxID, 10)},
		"count":        {strconv.Itoa(count)},
		"page":         {strconv.Itoa(page)},
		"base_app":     {strconv.Itoa(baseApp)},
		"feature":      {strconv.Itoa(feature)},
		"trim_user":    {strconv.Itoa(trimUser)},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesUserTimeline NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesUserTimeline Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesUserTimeline ReadAll error")
	}
	r := &StatusesUserTimelineResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesUserTimeline Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesUserTimeline resp error:" + r.Error)
	}
	return r, nil
}
