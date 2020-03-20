// 各种结构体定义

package weibo

import (
	"net/http"
)

// StatusesShareResp 微博成功发送后的返回结构
type StatusesShareResp struct {
	IDStr string `json:"idstr"`
}

// Weibo 定义各种微博相关方法
type Weibo struct {
	client      *http.Client
	appkey      string
	appsecret   string
	redirecturi string
	username    string
	passwd      string
	userAgent   string
}

// MobileLoginResp 移动登录返回结构
type MobileLoginResp struct {
	Retcode int                    `json:"retcode"`
	Msg     string                 `json:"msg"`
	Data    map[string]interface{} `json:"data"`
}

// PreLoginResp PC端prelogin登录返回结构
type PreLoginResp struct {
	Retcode    int    `json:"retcode"`
	Servertime int    `json:"servertime"`
	Pcid       string `json:"pcid"`
	Nonce      string `json:"nonce"`
	Pubkey     string `json:"pubkey"`
	Rsakv      string `json:"rsakv"`
	IsOpenlock int    `json:"is_openlock"`
	Showpin    int    `json:"showpin"`
	Exectime   int    `json:"exectime"`
}

// SsoLoginResp PC端ssologin登录返回结构
type SsoLoginResp struct {
	Retcode            string   `json:"retcode"`
	Ticket             string   `json:"ticket"`
	UID                string   `json:"uid"`
	Nick               string   `json:"nick"`
	CrossDomainURLList []string `json:"crossDomainUrlList"`
}

// RedirectResp 微博回调httpbin.org/get返回结构
type RedirectResp struct {
	Args map[string]string `json:"args"`
}

// TokenResp accesstoken返回结构
type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	UID         string `json:"uid"`
	IsRealName  string `json:"isRealName"`
}

// TokenInfoResp get_token_info接口返回结构
type TokenInfoResp struct {
	UID      string `json:"uid"`
	Appkey   string `json:"appkey"`
	Scope    string `json:"scope"`
	CreateAt string `json:"create_at"`
	ExpireIn string `json:"expire_in"`
}
