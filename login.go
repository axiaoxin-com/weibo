// 登录接口

package weibo

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

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
