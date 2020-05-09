// https://open.weibo.com/wiki/2/statuses/show
// 请求参数
// screen_name	false	string	需要查询的用户昵称。
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// id	true	int64	需要获取的微博ID。
//
// 注意事项
// 查询的微博必须是授权用户发出的，非授权用户发出的将不返回数据
//
// 返回字段说明
//   created_at	string	微博创建时间
//   id	int64	微博ID
//   mid	int64	微博MID
//   idstr	string	字符串型的微博ID
//   text	string	微博信息内容
//   source	string	微博来源
//   favorited	boolean	是否已收藏，true：是，false：否
//   truncated	boolean	是否被截断，true：是，false：否
//   in_reply_to_status_id	string	（暂未支持）回复ID
//   in_reply_to_user_id	string	（暂未支持）回复人UID
//   in_reply_to_screen_name	string	（暂未支持）回复人昵称
//   thumbnail_pic	string	缩略图片地址，没有时不返回此字段
//   bmiddle_pic	string	中等尺寸图片地址，没有时不返回此字段
//   original_pic	string	原始图片地址，没有时不返回此字段
//   geo	object	地理信息字段 详细
//   user	object	微博作者的用户信息字段 详细
//   retweeted_status	object	被转发的原微博信息字段，当该微博为转发微博时返回 详细
//   reposts_count	int	转发数
//   comments_count	int	评论数
//   attitudes_count	int	表态数
//   mlevel	int	暂未支持
//   visible	object	微博的可见性及指定可见分组信息。该object中type取值，0：普通微博，1：私密微博，3：指定分组微博，4：密友微博；list_id为分组的组号
//   pic_ids	object	微博配图ID。多图时返回多图ID，用来拼接图片url。用返回字段thumbnail_pic的地址配上该返回字段的图片ID，即可得到多个图片url。
//   ad	object array	微博流内的推广微博ID

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// StatusesShow 根据微博ID获取单条微博内容
func (w *Weibo) StatusesShow(token string, id int64) (*StatusesShowResp, error) {
	apiURL := "https://api.weibo.com/2/statuses/show.json"
	data := url.Values{
		"access_token": {token},
		"id":           {strconv.FormatInt(id, 10)},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShow NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShow Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShow ReadAll error")
	}
	r := &StatusesShowResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShow Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesShow resp error:" + r.Error)
	}
	return r, nil
}
