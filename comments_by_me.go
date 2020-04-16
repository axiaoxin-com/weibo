// https://open.weibo.com/wiki/2/comments/by_me
// 请求参数
//   access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
//   since_id	false	int64	若指定此参数，则返回ID比since_id大的评论（即比since_id时间晚的评论），默认为0。
//   max_id	false	int64	若指定此参数，则返回ID小于或等于max_id的评论，默认为0。
//   count	false	int	单页返回的记录条数，默认为50。
//   page	false	int	返回结果的页码，默认为1。
//   filter_by_source	false	int	来源筛选类型，0：全部、1：来自微博的评论、2：来自微群的评论，默认为0。
// 返回字段说明
//   created_at	string	评论创建时间
//   id	int64	评论的ID
//   text	string	评论的内容
//   source	string	评论的来源
//   user	object	评论作者的用户信息字段 详细
//   mid	string	评论的MID
//   idstr	string	字符串型的评论ID
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

// CommentsByMe 获取当前登录用户所发出的评论列表
// sinceID 返回ID比since_id大的评论（即比since_id时间晚的评论）
// maxID 返回ID小于或等于max_id的评论
// count 单页返回的记录条数
// page 返回结果的页码
// filterBySource 来源筛选类型，0：全部、1：来自微博的评论、2：来自微群的评论
func (w *Weibo) CommentsByMe(token string, sinceID, maxID, count, page, filterBySource int) (*CommentsByMeResp, error) {
	apiURL := "https://api.weibo.com/2/comments/by_me.json"
	data := url.Values{
		"access_token":     {token},
		"since_id":         {strconv.Itoa(sinceID)},
		"max_id":           {strconv.Itoa(maxID)},
		"count":            {strconv.Itoa(count)},
		"page":             {strconv.Itoa(page)},
		"filter_by_source": {strconv.Itoa(filterBySource)},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsByMe NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsByMe Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsByMe ReadAll error")
	}
	r := &CommentsByMeResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo CommentsByMe Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo CommentsByMe resp error:" + r.Error)
	}
	return r, nil
}
