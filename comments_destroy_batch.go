// https://open.weibo.com/wiki/2/comments/destroy_batch
// 请求参数
//   access_token	true	string	采用 OAuth 授权方式为必填参数， OAuth 授权后获得。
//   cids	true	int64	需要删除的评论 ID ，用半角逗号隔开，最多 20 个。

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

// CommentsDestroyBatch 根据评论 ID 批量删除评论
// cids 需要删除的评论 ID
func (w *Weibo) CommentsDestroyBatch(token string, cids ...int64) (*CommentsDestroyBatchResp, error) {
	apiURL := "https://api.weibo.com/2/comments/destroy_batch.json"
	sCids := []string{}
	for _, cid := range cids {
		sCids = append(sCids, strconv.FormatInt(cid, 10))
	}
	cidsStr := strings.Join(sCids, ",")
	data := url.Values{
		"access_token": {token},
		"cids":         {cidsStr},
	}
	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsDestroyBatch NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsDestroyBatch Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsDestroyBatch ReadAll error")
	}
	e := &ErrorResp{}
	json.Unmarshal(body, e)
	if e.Error != "" && e.ErrorCode != 0 {
		return nil, errors.New("weibo CommentsDestroyBatch resp error:" + e.Error)
	}
	r := &CommentsDestroyBatchResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo CommentsDestroyBatch Unmarshal error:"+string(body))
	}
	return r, nil
}
