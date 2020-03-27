package weibo_test

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin-com/chaojiying"
	"github.com/axiaoxin-com/weibo"
)

// 注册验证码破解函数示例
func Example_loginWithCrackFunc() {
	// 登录遇到验证码时
	// 如果有注册你自己的破解函数则会尝试使用你注册的函数进行验证码破解
	// 破解失败则采用默认的人工手动处理的方式手工输入保存在临时目录中的weibo_pin.png中的验证码

	// 从环境变量获取必须的账号信息
	appkey := os.Getenv("weibo_app_key")
	appsecret := os.Getenv("weibo_app_secret")
	username := os.Getenv("weibo_username")
	passwd := os.Getenv("weibo_passwd")
	redirecturi := os.Getenv("weibo_redirect_uri")
	securitydomain := os.Getenv("weibo_security_domain")

	// 初始化客户端
	weibo := weibo.New(appkey, appsecret, username, passwd, redirecturi)

	// 使用超级鹰破解验证码
	// 初始化超级鹰客户端
	chaojiyingUser := os.Getenv("chaojiying_user")
	chaojiyingPass := os.Getenv("chaojiying_pass")
	chaojiyingAccount := chaojiying.Account{User: chaojiyingUser, Pass: chaojiyingPass}
	cracker, err := chaojiying.New([]chaojiying.Account{chaojiyingAccount})
	if err != nil {
		log.Println(err)
	}
	// 将破解函数注册到微博客户端
	// 破解函数的声明为 func(io.Reader) (string, error)，只要符合此签名的函数就可以注册
	// RegisterCrackPinFunc 可以传入多个破解函数，会逐个尝试
	// 这里的Cr4ck即为chaojiying中的破解函数
	weibo.RegisterCrackPinFunc(cracker.Cr4ck)
	fmt.Println("验证码破解方法注册成功")

	// 登录微博 遇到验证码将自动识别
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
	status := fmt.Sprintf("文本内容 http://%s", securitydomain)
	resp, err := weibo.StatusesShare(token.AccessToken, status, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("微博发送成功 详情点击 http://weibo.com/" + resp.User.ProfileURL)
}
