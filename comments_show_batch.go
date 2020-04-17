// https://open.weibo.com/wiki/2/comments/show_batch
// 请求参数
//   access_token	true	string	采用 OAuth 授权方式为必填参数， OAuth 授权后获得。
//   cids	true	int64	需要查询的批量评论 ID ，用半角逗号分隔，最大 50

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

// CommentsShowBatch 根据评论 ID 批量返回评论信息
// cids 需要查询的批量评论 ID
func (w *Weibo) CommentsShowBatch(token string, cids ...int64) (*CommentsShowBatchResp, error) {
	apiURL := "https://api.weibo.com/2/comments/show_batch.json"
	sCids := []string{}
	for _, cid := range cids {
		sCids = append(sCids, strconv.FormatInt(cid, 10))
	}
	cidsStr := strings.Join(sCids, ",")
	data := url.Values{
		"access_token": {token},
		"cids":         {cidsStr},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsShowBatch NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsShowBatch Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo CommentsShowBatch ReadAll error")
	}
	e := &ErrorResp{}
	json.Unmarshal(body, e)
	if e.Error != "" && e.ErrorCode != 0 {
		return nil, errors.New("weibo CommentsShowBatch resp error:" + e.Error)
	}
	r := &CommentsShowBatchResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo CommentsShowBatch Unmarshal error:"+string(body))
	}
	return r, nil
}
