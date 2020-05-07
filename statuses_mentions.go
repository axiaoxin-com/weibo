// https://open.weibo.com/wiki/2/statuses/mentions
// 只返回授权用户的微博，非授权用户的微博将不返回；
// 请求参数
//   access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
//   since_id	false	int64	若指定此参数，则返回ID比since_id大的评论（即比since_id时间晚的评论），默认为0。
//   max_id	false	int64	若指定此参数，则返回ID小于或等于max_id的评论，默认为0。
//   count	false	int	单页返回的记录条数，默认为50。
//   page	false	int	返回结果的页码，默认为1。
//   filter_by_source	false	int	来源筛选类型，0：全部、1：来自微博的评论、2：来自微群的评论，默认为0。
//   filter_by_author	false	int	作者筛选类型，0：全部、1：我关注的人、2：陌生人，默认为0。
//   filter_by_type	false	int	原创筛选类型，0：全部微博、1：原创的微博，默认为0。
//
// 返回字段说明
// created_at	string	微博创建时间
// id	int64	微博ID
// mid	int64	微博MID
// idstr	string	字符串型的微博ID
// text	string	微博信息内容
// source	string	微博来源
// favorited	boolean	是否已收藏，true：是，false：否
// truncated	boolean	是否被截断，true：是，false：否
// in_reply_to_status_id	string	（暂未支持）回复ID
// in_reply_to_user_id	string	（暂未支持）回复人UID
// in_reply_to_screen_name	string	（暂未支持）回复人昵称
// thumbnail_pic	string	缩略图片地址，没有时不返回此字段
// bmiddle_pic	string	中等尺寸图片地址，没有时不返回此字段
// original_pic	string	原始图片地址，没有时不返回此字段
// geo	object	地理信息字段 详细
// user	object	微博作者的用户信息字段 详细
// retweeted_status	object	被转发的原微博信息字段，当该微博为转发微博时返回 详细
// reposts_count	int	转发数
// comments_count	int	评论数
// attitudes_count	int	表态数
// mlevel	int	暂未支持
// visible	object	微博的可见性及指定可见分组信息。该object中type取值，0：普通微博，1：私密微博，3：指定分组微博，4：密友微博；list_id为分组的组号
// pic_ids	object	微博配图ID。多图时返回多图ID，用来拼接图片url。用返回字段thumbnail_pic的地址配上该返回字段的图片ID，即可得到多个图片url。
// ad	object array	微博流内的推广微博ID

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// StatusesMentions 获取最新的提到登录用户的微博列表，即@我的微博
// sinceID 返回ID比since_id大的评论（即比since_id时间晚的评论）
// maxID 返回ID小于或等于max_id的评论
// count 单页返回的记录条数
// page 返回结果的页码
// filterBySource 来源筛选类型，0：全部、1：来自微博的评论、2：来自微群的评论
// filterByAuthor 作者筛选类型，0：全部、1：我关注的人、2：陌生人
// filterByType 原创筛选类型，0：全部微博、1：原创的微博
func (w *Weibo) StatusesMentions(token string, sinceID, maxID int64, count, page, filterBySource, filterByAuthor, filterByType int) (*StatusesMentionsResp, error) {
	apiURL := "https://api.weibo.com/2/statuses/mentions.json"
	data := url.Values{
		"access_token":     {token},
		"since_id":         {strconv.FormatInt(sinceID, 10)},
		"max_id":           {strconv.FormatInt(maxID, 10)},
		"count":            {strconv.Itoa(count)},
		"page":             {strconv.Itoa(page)},
		"filter_by_source": {strconv.Itoa(filterBySource)},
		"filter_by_author": {strconv.Itoa(filterByAuthor)},
		"filter_by_type":   {strconv.Itoa(filterByType)},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesMentions NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesMentions Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesMentions ReadAll error")
	}
	r := &StatusesMentionsResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesMentions Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesMentions resp error:" + r.Error)
	}
	return r, nil
}
