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

// RespStatusesCount StatusesCount 接口返回结构
type RespStatusesCount []struct {
	ID                    int64  `json:"id"`
	Idstr                 string `json:"idstr"`
	Comments              int    `json:"comments"`
	Reposts               int    `json:"reposts"`
	Attitudes             int    `json:"attitudes"`
	NumberDisplayStrategy struct {
		ApplyScenarioFlag    int    `json:"apply_scenario_flag"`
		DisplayTextMinNumber int    `json:"display_text_min_number"`
		DisplayText          string `json:"display_text"`
	} `json:"number_display_strategy"`
}

// StatusesCount 批量获取指定微博的转发数评论数
// ids 需要获取数据的微博ID，最多不超过100个。
func (w *Weibo) StatusesCount(token string, ids ...int64) (*RespStatusesCount, error) {
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
	e := &RespError{}
	json.Unmarshal(body, e)
	if e.Error != "" && e.ErrorCode != 0 {
		return nil, errors.New("weibo StatusesCount resp error:" + e.Error)
	}
	r := &RespStatusesCount{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesCount Unmarshal error:"+string(body))
	}
	return r, nil
}
