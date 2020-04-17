// https://open.weibo.com/wiki/2/comments/reply
// 请求参数
//   access_token	true	string	采用 OAuth 授权方式为必填参数， OAuth 授权后获得。
//   cid	true	int64	需要回复的评论 ID 。
//   id	true	int64	需要评论的微博 ID 。
//   comment	true	string	回复评论内容，必须做 URLencode ，内容不超过 140 个汉字。
//   without_mention	false	int	回复中是否自动加入“回复@用户名”， 0 ：是、 1 ：否，默认为 0 。
//   comment_ori	false	int	当评论转发微博时，是否评论给原微博， 0 ：否、 1 ：是，默认为 0 。
//   rip	false	string	开发者上报的操作用户真实 IP ，形如： 211.156.0.1 。

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

// CommentsReply 回复一条评论
// cid	需要回复的评论 ID 。
// id	需要评论的微博 ID 。
// comment 回复评论内容，内容不超过 140 个汉字。
// withoutMention	回复中是否自动加入“回复@用户名”， 0 ：是、 1 ：否。
// commentOri	当评论转发微博时，是否评论给原微博， 0 ：否、 1 ：是。
func (w *Weibo) CommentsReply(token string, cid, id int64, comment string, withoutMention, commentOri int) (*CommentsReplyResp, error) {
	apiURL := "https://api.weibo.com/2/comments/reply.json"
	data := url.Values{
		"access_token":    {token},
		"cid":             {strconv.FormatInt(cid, 10)},
		"id":              {strconv.FormatInt(id, 10)},
		"comment":         {comment},
		"without_mention": {strconv.Itoa(withoutMention)},
		"comment_ori":     {strconv.Itoa(commentOri)},
		"rip":             {RealIP()},
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsReply NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsReply Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsReply ReadAll error")
	}
	r := &CommentsReplyResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo CommentsReply Unmarshal error:"+string(body))
	}
	if r.Error != "" && r.ErrorCode != 0 {
		return nil, errors.New("weibo CommentsReply resp error:" + r.Error)
	}
	return r, nil
}
