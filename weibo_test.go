package weibo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

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
