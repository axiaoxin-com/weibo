package weibo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// Bug(testing_with_pin): Go在运行单元测试的时候，被测试的方法里面如果有fmt.Scan这种等待接收外部输入的操作的话它们都不会发生等待，而是会直接跳过，要是在一个死循环里面等待接收输入的方法跑测试就会陷进死循环里，原因是跑测试的时候会把输入重定向到/dev/null
// 跑测试时如果登录触发验证码时会导致无法进行验证码的输入，无法进行后续测试
func TestStatusesShare(t *testing.T) {
	appkey := os.Getenv("weibo_app_key")
	appsecret := os.Getenv("weibo_app_secret")
	username := os.Getenv("weibo_username")
	passwd := os.Getenv("weibo_passwd")
	redirecturi := os.Getenv("weibo_redirect_uri")
	securitydomain := os.Getenv("weibo_security_domain")
	weibo := New(appkey, appsecret, username, passwd, redirecturi)
	t.Log("PCLogin...")
	if err := weibo.PCLogin(); err != nil {
		t.Fatal(err)
	}

	t.Log("Authorize")
	code, err := weibo.Authorize()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("AccessToken")
	token, err := weibo.AccessToken(code)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("StatusesShare text")
	status := fmt.Sprintf("unit test at %s http://%s", time.Now().Format("2006-01-02 15:04:05"), securitydomain)
	resp, err := weibo.StatusesShare(token.AccessToken, status, nil)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(2 * time.Second)
	t.Log("StatusesShare pic")
	pic, err := os.Open("./example/pic.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer pic.Close()
	resp, err = weibo.StatusesShare(token.AccessToken, status, pic)
	if err != nil {
		t.Error(err)
	}
	t.Log("http://weibo.com/" + resp.User.ProfileURL)
}
