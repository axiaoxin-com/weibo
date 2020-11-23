// Package weibo 新浪微博 SDK
package weibo

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// RespError 微博接口的错误结果返回结构
type RespError struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Request   string `json:"request"`
}

// HTTPTimeout 请求超时时间 默认 10 秒
var HTTPTimeout time.Duration = time.Second * 10

// Weibo 实例，在其上实现各类接口
type Weibo struct {
	client        *http.Client
	appkey        string
	appsecret     string
	redirecturi   string
	username      string
	passwd        string
	userAgent     string
	crackPinFuncs []CrackPinFunc
	token         string
}

// New 创建Weibo实例
// appkey 微博开放平台应用的 appkey
// appsecret 微博开放平台应用的 appsecret
// username 需要发微博的微博登录账号，用于模拟登录直接获取授权码
// password 需要发微博的微博登录密码，用于模拟登录直接获取授权码
// redirecturi 微博开发平台应用的回调 URL
func New(appkey, appsecret, username, passwd, redirecturi string) *Weibo {
	// 设置cookiejar后续请求会自动带cookie保持会话
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: HTTPTimeout,
	}
	return &Weibo{
		client:      client,
		appkey:      appkey,
		appsecret:   appsecret,
		redirecturi: redirecturi,
		username:    username,
		passwd:      passwd,
		userAgent:   RandUA(),
	}
}
