// https://open.weibo.com/wiki/Oauth2/access_token
// 请求参数：
//    client_id	true	string	申请应用时分配的AppKey。
//    client_secret	true	string	申请应用时分配的AppSecret。
//    grant_type	true	string	请求的类型，填写authorization_code
// grant_type为authorization_code时:
//    code	true	string	调用authorize获得的code值。
//    redirect_uri	true	string	回调地址，需需与注册应用里的回调地址一致。
// 返回数据：
//    access_token	string	用户授权的唯一票据，用于调用微博的开放接口，同时也是第三方应用验证微博用户登录的唯一票据，第三方应用应该用该票据和自己应用内的用户建立唯一影射关系，来识别登录状态，不能使用本返回值里的UID字段来做登录识别。
//    expires_in	string	access_token的生命周期，单位是秒数。
//    remind_in	string	access_token的生命周期（该参数即将废弃，开发者请使用expires_in）。
//    uid	string	授权用户的UID，本字段只是为了方便开发者，减少一次user/show接口调用而返回的，第三方应用不能用此字段作为用户登录状态的识别，只有access_token才是用户授权的唯一票据。

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// AccessToken 传入授权码请求access_token接口，返回token对象
func (w *Weibo) AccessToken(code string) (*TokenResp, error) {
	tokenURL := "https://api.weibo.com/oauth2/access_token"
	data := url.Values{
		"client_id":     {w.appkey},
		"client_secret": {w.appsecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {w.redirecturi},
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo AccessToken NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo AccessToken Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo AccessToken ReadAll error")
	}
	tokenResp := &TokenResp{}
	if err := json.Unmarshal(body, tokenResp); err != nil {
		return nil, errors.Wrap(err, "weibo AccessToken Unmarshal error")
	}
	if tokenResp.AccessToken == "" {
		return nil, errors.New("weibo AccessToken get token failed." + string(body))
	}
	return tokenResp, nil
}
