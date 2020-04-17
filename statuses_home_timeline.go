// https://open.weibo.com/wiki/2/statuses/home_timeline
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// since_id	false	int64	若指定此参数，则返回ID比since_id大的微博（即比since_id时间晚的微博），默认为0。
// max_id	false	int64	若指定此参数，则返回ID小于或等于max_id的微博，默认为0。
// count	false	int	单页返回的记录条数，最大不超过100，默认为20。
// page	false	int	返回结果的页码，默认为1。
// base_app	false	int	是否只获取当前应用的数据。0为否（所有数据），1为是（仅当前应用），默认为0。
// feature	false	int	过滤类型ID，0：全部、1：原创、2：图片、3：视频、4：音乐，默认为0。
// trim_user	false	int	返回值中user字段开关，0：返回完整user字段、1：user字段仅返回user_id，默认为0。
package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// StatusesHomeTimeline 获取当前登录用户及其所关注（授权）用户的最新微博
// sinceID 返回ID比sinceID大的微博（即比since_id时间晚的微博）
// maxID 返回ID小于或等于max_id的微博0。
// count 单页返回的记录条数，最大不超过100。
// page	返回结果的页码。
// baseApp 是否只获取当前应用的数据。0为否（所有数据），1为是（仅当前应用）
// feature 过滤类型ID，0：全部、1：原创、2：图片、3：视频、4：音乐0。
// trim_user 返回值中user字段开关，0：返回完整user字段、1：user字段仅返回user_id。
func (w *Weibo) StatusesHomeTimeline(token string, sinceID, maxID int64, count, page, baseApp, feature, trimUser int) (*StatusesHomeTimelineResp, error) {
	apiURL := "https://api.weibo.com/2/statuses/home_timeline.json"
	data := url.Values{
		"access_token": {token},
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
		return nil, errors.Wrap(err, "weibo StatusesHomeTimeline NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesHomeTimeline Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesHomeTimeline ReadAll error")
	}
	r := &StatusesHomeTimelineResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesHomeTimeline Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesHomeTimeline resp error:" + r.Error)
	}
	return r, nil
}
