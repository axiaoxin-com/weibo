// https://open.weibo.com/wiki/Oauth2/get_token_info
// 请求参数：
//    access_token：用户授权时生成的access_token。
// 返回数据：
//    uid	string	授权用户的uid。
//    appkey	string	access_token所属的应用appkey。
//    scope	string	用户授权的scope权限。
//    create_at	string	access_token的创建时间，从1970年到创建时间的秒数。
//    expire_in	string	access_token的剩余时间，单位是秒数。

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// TokenInfo 获取用户 access_token 的授权相关信息，包括授权时间，过期时间和 scope 权限
func (w *Weibo) TokenInfo(token string) (*TokenInfoResp, error) {
	apiURL := "https://api.weibo.com/oauth2/get_token_info"
	data := url.Values{
		"access_token": {token},
	}
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo TokenInfo NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo TokenInfo Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo TokenInfo ReadAll error")
	}
	tokenInfoResp := &TokenInfoResp{}
	if err := json.Unmarshal(body, tokenInfoResp); err != nil {
		return nil, errors.Wrap(err, "weibo TokenInfo Unmarshal error")
	}
	return tokenInfoResp, nil
}
