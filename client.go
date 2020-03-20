// Package weibo 封装微博API
package weibo

import (
	"net/http"
	"net/http/cookiejar"
)

// New 创建Weibo实例
// appkey 微博开放平台appkey
// appsecret 微博开放平台appsecret
// username 微博登录账号
// password 微博密码
// redirecturi 微博开发平台app设置的回调url
func New(appkey, appsecret, username, passwd, redirecturi string) *Weibo {
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
