# weibo

新浪微博 Golang 版 SDK

微博API: <https://open.weibo.com/wiki/微博API>

## 功能

代码组织结构已按新浪微博提供的接口拆分，已支持的功能列表：

- [模拟电脑浏览器登录](https://github.com/axiaoxin-com/weibo/blob/master/login.go#L200)
- [注册验证码破解函数](https://github.com/axiaoxin-com/weibo/blob/master/login.go#L81)
- [模拟移动端登录](https://github.com/axiaoxin-com/weibo/blob/master/login.go#L30)
- [获取 Authorize Code](https://github.com/axiaoxin-com/weibo/blob/master/authorize.go)
- [获取 Access Token](https://github.com/axiaoxin-com/weibo/blob/master/access_token.go)
- [查询 Access Token 信息](https://github.com/axiaoxin-com/weibo/blob/master/get_token_info.go)
- [分享一条链接到微博（发微博）](https://github.com/axiaoxin-com/weibo/blob/master/statuses_share.go)

## 亮点

#### 模拟微博登录自动获取授权码并取得token

使用账号密码模拟登录微博后获取授权码，从url中取得授权码后再获取token，过程中无需人工干预。

#### 支持登录时验证码识别

默认触发验证码时，将验证码保存在临时目录中，提示用户人工处理，手动输入验证码后继续后续逻辑，期间会尝试显示验证码图片，若失败则需要人工去提示路径下打开图片。

支持注册破解验证码的函数，注册后触发验证码时，优先使用注册的函数识别验证码，如果识别失败则仍然采用提示用户手动输入。

破解函数的声明为 `func(io.Reader) (string, error)` ，只要符合此签名的函数就可以调用 `RegisterCrackPinFunc` 方法注册。`RegisterCrackPinFunc` 可以传入多个破解函数，会逐个尝试。

## 安装

```
go get -u -v github.com/axiaoxin-com/weibo
```

## 在线文档

<https://godoc.org/github.com/axiaoxin-com/weibo>

## 使用示例

### 发送纯文本内容的微博

[example/text.go](https://github.com/axiaoxin-com/weibo/blob/master/example/text.go)

```go
// 发送文本内容示例

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin-com/weibo"
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

[example/text_pic.go](https://github.com/axiaoxin-com/weibo/blob/master/example/text_pic.go)

```go
// 发送带图片的文本内容示例

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axiaoxin-com/weibo"
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

[example/crackfunc.go](https://github.com/axiaoxin-com/weibo/blob/master/example/crackfunc.go)

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

	"github.com/axiaoxin-com/chaojiying"
	"github.com/axiaoxin-com/weibo"
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

## 难点

#### 模拟登录

要想代码层面直接获取到授权码，必须要在微博应用的授权页面进行模拟浏览器登录。
登录参数会被js代码处理，需要翻译对应的 js 代码为 go ，crypto 包不熟，这里花了一些时间。

#### 如何处理模拟登录时遇到的验证码

最开始的方案是客户端实例化的时候可以接收授权码参数，如果直接传入了授权码则跳过登录逻辑直接去获取 token 就不会触发验证码。
但是这个接口使用体验太差，每次需要先去浏览器中手工授权获取 URL 中的授权码，再将其写入配置中运行时加载，过程相对麻烦；
而且如果为多个微博账号配置各自的授权码时极易忘记退出上一次登录的账号导致获取到的是其他账号的 token 。

在无法自己实现验证码自动识别的事实上，不如将验证码保存到本地，触发验证码时本地直接打开后手动输入即可，只需一次人工干预，
实际应用中通常也只有服务启动时需要登录授权操作，因此这样可以接受。

在调研后发现有超级鹰等类似的验证码识别平台，提供 http 接口破解验证码，于是实现了注册破解函数的模式。
用户可以实现自己的破解方法，只需按约定的出入参数实现即可，注册后当遇到验证码会自动调用处理，失败时再使用人工识别手动输入的方案兜底。


## 遇到的问题

<https://github.com/axiaoxin-com/weibo/issues?q=label%3Anote>
