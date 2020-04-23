// https://open.weibo.com/wiki/2/statuses/count
// 请求参数
//   access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// 	 ids 需要获取数据的微博ID，多个之间用逗号分隔，最多不超过100个。

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

// StatusesCount 批量获取指定微博的转发数评论数
// ids 需要获取数据的微博ID，最多不超过100个。
func (w *Weibo) StatusesCount(token string, ids ...int64) (*StatusesCountResp, error) {
	apiURL := "https://api.weibo.com/2/statuses/count.json"
	sIds := []string{}
	for _, id := range ids {
		sIds = append(sIds, strconv.FormatInt(id, 10))
	}
	idsStr := strings.Join(sIds, ",")
	data := url.Values{
		"access_token": {token},
		"ids":          {idsStr},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesCount NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesCount Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesCount ReadAll error")
	}
	e := &ErrorResp{}
	json.Unmarshal(body, e)
	if e.Error != "" && e.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesCount resp error:" + e.Error)
	}
	r := &StatusesCountResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesCount Unmarshal error:"+string(body))
	}
	return r, nil
}
