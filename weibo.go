// Package weibo 封装微博API
package weibo

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// NewWeibo 创建Weibo实例
// appkey 微博开放平台appkey
// appsecret 微博开放平台appsecret
// username 微博登录账号
// password 微博密码
// redirecturi 微博开发平台app设置的回调url
func NewWeibo(appkey, appsecret, username, passwd, redirecturi string) *Weibo {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	return &Weibo{
		client:      client,
		appkey:      appkey,
		appsecret:   appsecret,
		redirecturi: redirecturi,
		username:    username,
		passwd:      passwd,
		userAgent:   randUA(),
	}
}

// MobileLogin 移动端登录微博
// 该登录无法获取开放平台token
func (w *Weibo) MobileLogin() error {
	data := url.Values{
		"username":     {w.username},
		"password":     {w.passwd},
		"savestate":    {"1"},
		"r":            {"https://weibo.cn/"},
		"ec":           {"0"},
		"pagerefer":    {"https://weibo.cn/pub/"},
		"entry":        {"mweibo"},
		"wentry":       {""},
		"loginfrom":    {""},
		"client_id":    {""},
		"code":         {""},
		"qq":           {""},
		"mainpageflag": {"1"},
		"hff":          {""},
		"hfp":          {""},
	}
	logingURL := "https://passport.weibo.cn/sso/login"
	req, err := http.NewRequest("POST", logingURL, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin NewRequest error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "deflate, br") // no gzip
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://passport.weibo.cn")
	req.Header.Set("Referer", "https://passport.weibo.cn/signin/login?entry=mweibo&r=https%3A%2F%2Fweibo.cn%2F&backTitle=%CE%A2%B2%A9&vt=")
	resp, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin Do error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin ReadAll error")
	}
	loginResp := &MobileLoginResp{}
	if err := json.Unmarshal(body, loginResp); err != nil {
		return errors.Wrap(err, "weibo MobileLogin Unmarshal error")
	}
	if loginResp.Retcode != 20000000 {
		return errors.New("weibo MobileLogin loginResp Retcode error:" + string(body))
	}
	return nil
}

// PCLogin 电脑web登录
// 登录后才能成功获取开放平台授权码和token
func (w *Weibo) PCLogin() error {
	preloginURL := "https://login.sina.com.cn/sso/prelogin.php?"
	ssologinURL := "https://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.19)"
	// pinURL := "https://login.sina.com.cn/cgi/pin.php" // 登录验证码相关url

	// 请求prelogin 获得 servertime, nonce, pubkey, rsakv
	// 对账号进行base64编码 对应javascript中encodeURIComponent然后base64编码
	su := base64.StdEncoding.EncodeToString([]byte(w.username))
	req, err := http.NewRequest("GET", preloginURL, nil)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin NewRequest prelogin error")
	}
	params := url.Values{
		"entry":    {"weibo"},
		"su":       {su},
		"rsakt":    {"mod"},
		"checkpin": {"1"},
		"client":   {"ssologin.js(v1.4.19)"},
		"_":        {strconv.FormatInt(time.Now().UnixNano()/1e6, 10)},
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin Do prelogin error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin ReadAll prelogin error")
	}
	preLoginResp := &PreLoginResp{}
	if err := json.Unmarshal(body, preLoginResp); err != nil {
		return errors.Wrap(err, "weibo PCLogin Unmarshal preLoginResp error")
	}
	if preLoginResp.Retcode != 0 {
		return errors.New("weibo PCLogin preLoginResp Retcode error:" + string(body))
	}

	// 请求ssologin
	// 拼接明文js加密文件中得到
	encMsg := []byte(fmt.Sprint(preLoginResp.Servertime, "\t", preLoginResp.Nonce, "\n", w.passwd))
	// 创建公钥
	n, _ := new(big.Int).SetString(preLoginResp.Pubkey, 16)
	e, _ := new(big.Int).SetString("10001", 16)
	pubkey := &rsa.PublicKey{N: n, E: int(e.Int64())}
	// 加密公钥
	sp, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey, encMsg)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin EncryptPKCS1v15 error")
	}
	// 将加密信息转换为16进制
	hexsp := hex.EncodeToString([]byte(sp))
	data := url.Values{
		"entry":      {"account"},
		"gateway":    {"1"},
		"from":       {""},
		"savestate":  {"30"},
		"useticket":  {"1"},
		"pagerefer":  {""},
		"vsnf":       {"1"},
		"su":         {su},
		"service":    {"account"},
		"servertime": {fmt.Sprint(preLoginResp.Servertime, randInt(1, 20))},
		"nonce":      {preLoginResp.Nonce},
		"pwencode":   {"rsa2"},
		"rsakv":      {preLoginResp.Rsakv},
		"sp":         {hexsp},
		"sr":         {"1536 * 864"},
		"encoding":   {"UTF - 8"},
		"cdult":      {"3"},
		"domain":     {"sina.com.cn"},
		"prelt":      {"95"},
		"returntype": {"TEXT"},
	}
	req, err = http.NewRequest("POST", ssologinURL, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin NewRequest ssologin error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin Do ssologin error")
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo PCLogin ReadAll ssologin error")
	}
	ssoLoginResp := &SsoLoginResp{}
	if err := json.Unmarshal(body, ssoLoginResp); err != nil {
		return errors.Wrap(err, "weibo PCLogin Unmarshal ssoLoginResp error")
	}
	if ssoLoginResp.Retcode != "0" {
		return errors.New("weibo PCLogin ssoLoginResp Retcode error:" + string(body))
	}
	return w.loginSucceed(ssoLoginResp)
}

// loginSucceed 检查PCLogin后是否成功登录微博
func (w *Weibo) loginSucceed(resp *SsoLoginResp) error {
	// 请求login_url和home_url, 进一步验证登录是否成功
	s := strings.Split(strings.Split(resp.CrossDomainURLList[0], "ticket=")[1], "&ssosavestate=")
	loginURL := fmt.Sprintf("https://passport.weibo.com/wbsso/login?ticket=%s&ssosavestate=%s&callback=sinaSSOController.doCrossDomainCallBack&scriptId=ssoscript0&client=ssologin.js(v1.4.19)&_=%s", s[0], s[1], strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	req, err := http.NewRequest("GET", loginURL, nil)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed NewRequest loginURL error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	res, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed Do loginURL error")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed ReadAll loginURL error")
	}
	reg := regexp.MustCompile(`"uniqueid":"(.*?)"`)
	result := reg.FindAllStringSubmatch(string(body), -1)
	if len(result) == 0 {
		return errors.New("weibo loginSucceed uniqueid pattern not match")
	}
	uid := result[0][1]
	homeURL := fmt.Sprintf("https://weibo.com/u/%s/home", uid)
	req, err = http.NewRequest("GET", homeURL, nil)
	req.Header.Set("User-Agent", w.userAgent)
	res, err = w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed Do homeURL error")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed ReadAll homeURL error")
	}
	if !strings.Contains(string(body), "我的首页") {
		return errors.New("weibo loginSucceed login failed")
	}
	return nil
}

// AuthCode 请求微博授权，返回授权码
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
func (w *Weibo) AuthCode() (string, error) {
	authURL := "https://api.weibo.com/oauth2/authorize"
	referer := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", authURL, w.appkey, w.redirecturi)
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
		return "", errors.Wrap(err, "weibo AuthCode NewRequest error")
	}
	req.Header.Set("Referer", referer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "weibo AuthCode Do error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "weibo AuthCode ReadAll error")
	}
	redirectResp := &RedirectResp{}
	if err := json.Unmarshal(body, redirectResp); err != nil {
		return "", errors.Wrap(err, "weibo AuthCode Unmarshal error")
	}
	return redirectResp.Args["code"], nil
}

// AccessToken 请求access_token接口，返回token对象
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
	return tokenResp, nil
}

// StatusesShare 第三方分享一条链接到微博
// https://open.weibo.com/wiki/2/statuses/share
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// status	true	string	用户分享到微博的文本内容，必须做URLencode，内容不超过140个汉字，文本中不能包含“#话题词#”，同时文本中必须包含至少一个第三方分享到微博的网页URL，且该URL只能是该第三方（调用方）绑定域下的URL链接，绑定域在“我的应用 － 应用信息 － 基本应用信息编辑 － 安全域名”里设置。
// pic	false	binary	用户想要分享到微博的图片，仅支持JPEG、GIF、PNG图片，上传图片大小限制为<5M。上传图片时，POST方式提交请求，需要采用multipart/form-data编码方式。
// rip	false	string	开发者上报的操作用户真实IP，形如：211.156.0.1。
func (w *Weibo) StatusesShare(token, status string, pic io.Reader) error {
	apiURL := "https://api.weibo.com/2/statuses/share.json"
	ip := realip()
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	if pic == nil {
		data := url.Values{
			"access_token": {token},
			"status":       {status},
			"rip":          {ip},
		}
		bodyBuf = bytes.NewBufferString(data.Encode())
	} else {
		picWriter, err := writer.CreateFormFile("pic", "picname.png")
		if err != nil {
			return errors.Wrap(err, "weibo StatusesShare CreateFormFile error")
		}
		if _, err := io.Copy(picWriter, pic); err != nil {
			return errors.Wrap(err, "weibo StatusesShare io.Copy error")
		}

		if err := writer.WriteField("access_token", token); err != nil {
			return errors.Wrap(err, "weibo StatusesShare WriteField access_token error")
		}
		if err := writer.WriteField("status", status); err != nil {
			return errors.Wrap(err, "weibo StatusesShare WriteField status error")
		}
		if err := writer.WriteField("rip", ip); err != nil {
			return errors.Wrap(err, "weibo StatusesShare WriteField rip error")
		}
		writer.Close() // must close before new request
	}
	req, err := http.NewRequest("POST", apiURL, bodyBuf)
	if err != nil {
		return errors.Wrap(err, "weibo StatusesShare NewRequest error")
	}
	if pic == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo StatusesShare Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo StatusesShare ReadAll error")
	}
	sr := &StatusesShareResp{}
	if err := json.Unmarshal(body, sr); err != nil {
		return errors.Wrap(err, "weibo StatusesShare Unmarshal error:"+string(body))
	}
	if sr.IDStr == "" {
		return errors.New(string(body))
	}
	return nil
}

// TokenInfo 查询用户access_token的授权相关信息，包括授权时间，过期时间和scope权限
// https://open.weibo.com/wiki/Oauth2/get_token_info
// 请求参数：
//    access_token：用户授权时生成的access_token。
// 返回数据：
//    uid	string	授权用户的uid。
//    appkey	string	access_token所属的应用appkey。
//    scope	string	用户授权的scope权限。
//    create_at	string	access_token的创建时间，从1970年到创建时间的秒数。
//    expire_in	string	access_token的剩余时间，单位是秒数。
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
