// 发送带图片的文本内容示例

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin/weibo"
)

func main() {
	// 从环境变量获取必须的账号信息
	appkey := os.Getenv("weibo_app_key")
	appsecret := os.Getenv("weibo_app_secret")
	username := os.Getenv("weibo_username")
	passwd := os.Getenv("weibo_passwd")
	redirecturi := os.Getenv("weibo_redirect_uri")
	securitydomain := os.Getenv("weibo_security_domain")

	// 初始化客户端
	weibo := weibo.New(appkey, appsecret, username, passwd, redirecturi)

	// 登录微博
	if err := weibo.PCLogin(); err != nil {
		log.Fatal(err)
	}

	// 获取授权码
	code, err := weibo.Authorize()
	if err != nil {
		log.Fatal(err)
	}

	// 获取access token
	token, err := weibo.AccessToken(code)
	if err != nil {
		log.Fatal(err)
	}

	// 发送微博，必须带有安全域名链接
	status := fmt.Sprintf("文字带图片示例 http://%s", securitydomain)
	// 加载要发送的图片，加载方式只要是返回io.Reader都可以
	pic, err := os.Open("./pic.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer pic.Close()
	resp, err := weibo.StatusesShare(token.AccessToken, status, pic)
	if err != nil {
		log.Println(err)
	}
	log.Println("微博发送成功 详情点击 http://weibo.com/" + resp.User.ProfileURL)
}
