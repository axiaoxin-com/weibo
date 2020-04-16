// https://open.weibo.com/wiki/2/comments/create
// 请求参数
//   access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
//   comment	true	string	评论内容，必须做URLencode，内容不超过140个汉字。
//   id	true	int64	需要评论的微博ID。
//   comment_ori	false	int	当评论转发微博时，是否评论给原微博，0：否、1：是，默认为0。
//   rip	false	string	开发者上报的操作用户真实IP，形如：211.156.0.1。

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

// CommentsCreate 对一条微博进行评论
// comment	评论内容，不超过140个汉字。
// id 需要评论的微博ID。
// commentOri 当评论转发微博时，是否评论给原微博，0：否、1：是。
func (w *Weibo) CommentsCreate(token string, comment string, id int64, commentOri int) (*CommentsCreateResp, error) {
	apiURL := "https://api.weibo.com/2/comments/create.json"
	data := url.Values{
		"access_token": {token},
		"comment":      {comment},
		"id":           {strconv.FormatInt(id, 10)},
		"comment_ori":  {strconv.Itoa(commentOri)},
		"rip":          {RealIP()},
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsCreate NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsCreate Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsCreate ReadAll error")
	}
	r := &CommentsCreateResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo CommentsCreate Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo CommentsCreate resp error:" + r.Error)
	}
	return r, nil
}
