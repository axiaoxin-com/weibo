// https://open.weibo.com/wiki/2/comments/to_me
// 请求参数
//   access_token	true	string	采用 OAuth 授权方式为必填参数， OAuth 授权后获得。
//   since_id	false	int64	若指定此参数，则返回 ID 比 since_id 大的评论（即比 since_id 时间晚的评论），默认为 0 。
//   max_id	false	int64	若指定此参数，则返回 ID 小于或等于 max_id 的评论，默认为 0 。
//   count	false	int	单页返回的记录条数，默认为 50 。
//   page	false	int	返回结果的页码，默认为 1 。
//   filter_by_author	false	int	作者筛选类型， 0 ：全部、 1 ：我关注的人、 2 ：陌生人，默认为 0 。
//   filter_by_source	false	int	来源筛选类型， 0 ：全部、 1 ：来自微博的评论、 2 ：来自微群的评论，默认为 0 。
// 返回字段说明
//   created_at	string	评论创建时间
//   id	int64	评论的 ID
//   text	string	评论的内容
//   source	string	评论的来源
//   user	object	评论作者的用户信息字段 详细
//   mid	string	评论的 MID
//   idstr	string	字符串型的评论 ID
//   status	object	评论的微博信息字段 详细
//   reply_comment	object	评论来源评论，当本评论属于对另一评论的回复时返回此字段

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// CommentsToMe 获取当前登录用户所接收到的评论列表，返回结果和 RespCommentsByMe 相同
// sinceID 返回 ID 比 since_id 大的评论（即比 since_id 时间晚的评论）
// maxID 返回 ID 小于或等于 max_id 的评论
// count 单页返回的记录条数
// page 返回结果的页码
// filterByAuthor	作者筛选类型， 0 ：全部、 1 ：我关注的人、 2 ：陌生人。
// filterBySource 来源筛选类型， 0 ：全部、 1 ：来自微博的评论、 2 ：来自微群的评论
func (w *Weibo) CommentsToMe(token string, sinceID, maxID int64, count, page, filterByAuthor, filterBySource int) (*RespCommentsByMe, error) {
	apiURL := "https://api.weibo.com/2/comments/to_me.json"
	data := url.Values{
		"access_token":     {token},
		"since_id":         {strconv.FormatInt(sinceID, 10)},
		"max_id":           {strconv.FormatInt(maxID, 10)},
		"count":            {strconv.Itoa(count)},
		"page":             {strconv.Itoa(page)},
		"filter_by_author": {strconv.Itoa(filterByAuthor)},
		"filter_by_source": {strconv.Itoa(filterBySource)},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsToMe NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsToMe Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsToMe ReadAll error")
	}
	// 返回结构相同
	r := &RespCommentsByMe{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo CommentsToMe Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo CommentsToMe resp error:" + r.Error)
	}
	return r, nil
}
