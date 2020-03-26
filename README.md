# weibo

新浪微博 Golang 版 SDK

微博API: <https://open.weibo.com/wiki/微博API>

## 功能

代码组织结构已按新浪微博提供的接口拆分，已支持的功能列表：

- [模拟电脑浏览器登录](https://github.com/axiaoxin/weibo/blob/master/login.go#L200)
- [注册验证码破解函数](https://github.com/axiaoxin/weibo/blob/master/login.go#L81)
- [模拟移动端登录](https://github.com/axiaoxin/weibo/blob/master/login.go#L30)
- [获取 Authorize Code](https://github.com/axiaoxin/weibo/blob/master/authorize.go)
- [获取 Access Token](https://github.com/axiaoxin/weibo/blob/master/access_token.go)
- [查询 Access Token 信息](https://github.com/axiaoxin/weibo/blob/master/get_token_info.go)
- [分享一条链接到微博（发微博）](https://github.com/axiaoxin/weibo/blob/master/statuses_share.go)

## 特性

### 模拟微博登录自动获取授权码并取得token

使用账号密码模拟登录微博后获取授权码，从url中取得授权码后再获取token，过程中无需人工干预。

### 支持登录时验证码识别

默认触发验证码时，将验证码保存在临时目录中，提示用户人工处理，手动输入验证码后继续后续逻辑，期间会尝试显示验证码图片，若失败则需要人工去提示路径下打开图片。

支持注册破解验证码的函数，注册后触发验证码时，优先使用注册的函数识别验证码，如果识别失败则仍然采用提示用户手动输入。

破解函数的声明为 `func(io.Reader) (string, error)` ，只要符合此签名的函数就可以调用 `RegisterCrackPinFunc` 方法注册。`RegisterCrackPinFunc` 可以传入多个破解函数，会逐个尝试。

## 安装

```
go get -u -v github.com/axiaoxin/weibo
```

## 使用示例

### 发送纯文本内容的微博

[example/text.go](https://github.com/axiaoxin/weibo/blob/master/example/text.go)

```go
// 发送文本内容示例

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
	status := fmt.Sprintf("文本内容 http://%s", securitydomain)
	resp, err := weibo.StatusesShare(token.AccessToken, status, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("微博发送成功 详情点击 http://weibo.com/" + resp.User.ProfileURL)
}
```

### 发送文字内容带图片的微博

[example/text_pic.go](https://github.com/axiaoxin/weibo/blob/master/example/text_pic.go)

```go
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
```

### 注册破解函数

[example/crackfunc.go](https://github.com/axiaoxin/weibo/blob/master/example/crackfunc.go)

```go
// 注册验证码破解函数示例
// 登录遇到验证码时
// 如果有注册你自己的破解函数则会尝试使用你注册的函数进行验证码破解
// 破解失败则采用默认的人工手动处理的方式手工输入保存在临时目录中的weibo_pin.png中的验证码

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin/chaojiying"
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
```
