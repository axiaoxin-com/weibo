// https://open.weibo.com/wiki/Oauth2/authorize
// 请求参数：
//    client_id	true	string	申请应用时分配的AppKey。
//    redirect_uri	true	string	授权回调地址，站外应用需与设置的回调地址一致，站内应用需填写canvas page的地址。
//    scope	false	string	申请scope权限所需参数，可一次申请多个scope权限，用逗号分隔。使用文档
//    state	false	string	用于保持请求和回调的状态，在回调时，会在Query Parameter中回传该参数。开发者可以用这个参数验证请求有效性，也可以记录用户请求授权页前的位置。这个参数可用于防止跨站请求伪造（CSRF）攻击
//    display	false	string	授权页面的终端类型，取值见下面的说明。
//    forcelogin	false	boolean	是否强制用户重新登录，true：是，false：否。默认false。
//    language	false	string	授权页语言，缺省为中文简体版，en为英文版。英文版测试中，开发者任何意见可反馈至 @微博API
// display说明：
//    default	默认的授权页面，适用于web浏览器。
//    mobile	移动终端的授权页面，适用于支持html5的手机。注：使用此版授权页请用 https://open.weibo.cn/oauth2/authorize 授权接口
//    wap	wap版授权页面，适用于非智能手机。
//    client	客户端版本授权页面，适用于PC桌面应用。
//    apponweibo	默认的站内应用授权页，授权后不返回access_token，只刷新站内应用父框架。
// 返回数据：
//    code	string	用于第二步调用oauth2/access_token接口，获取授权后的access token。
//    state	string	如果传递参数，会回传该参数。

package weibo

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Authorize 请求微博授权，返回授权码
func (w *Weibo) Authorize() (string, error) {
	authURL := "https://api.weibo.com/oauth2/authorize"
	referer := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", authURL, w.appkey, w.redirecturi)
	log.Println("[DEBUG] weibo Authorize referer:", referer)
	data := url.Values{
		"client_id":       {w.appkey},
		"response_type":   {"code"},
		"redirect_uri":    {w.redirecturi},
		"action":          {"submit"},
		"userId":          {w.username},
		"passwd":          {w.passwd},
		"isLoginSina":     {"0"},
		"from":            {""},
		"regCallback":     {""},
		"state":           {""},
		"ticket":          {""},
		"withOfficalFlag": {"0"},
	}
	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", errors.Wrap(err, "weibo Authorize NewRequest error")
	}
	req.Header.Set("Referer", referer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "weibo Authorize Do error")
	}
	defer resp.Body.Close()
	uri := resp.Request.URL.String()
	s := strings.Split(uri, "code=")
	if len(s) != 2 {
		return "", errors.New("weibo Authorize get code from uri error, authorize code url=" + uri)
	}
	return s[1], nil
}
