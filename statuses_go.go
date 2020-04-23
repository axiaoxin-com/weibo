// https://open.weibo.com/wiki/2/statuses/go
// source	false	string	采用 OAuth 授权方式不需要此参数，其他授权方式为必填参数，数值为应用的 AppKey 。
// access_token	false	string	采用 OAuth 授权方式为必填参数，其他授权方式不需要此参数， OAuth 授权后获得。
// uid	true	int64	需要跳转的用户 ID 。
// id	true	int64	需要跳转的微博 ID 。

package weibo

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// StatusesGo 根据 ID 返回对应微博页面跳转 URL
// uid	int64	需要跳转的用户 ID 。
// id	int64	需要跳转的微博 ID 。
func (w *Weibo) StatusesGo(token string, uid, id int64) (string, error) {
	apiURL := "https://api.weibo.com/2/statuses/go"
	data := url.Values{
		"access_token": {token},
		"source":       {w.appkey},
		"uid":          {strconv.FormatInt(uid, 10)},
		"id":           {strconv.FormatInt(id, 10)},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "weibo StatusesGoURL NewRequest error")
	}
	req.URL.RawQuery = data.Encode()
	return req.URL.String(), nil
}
