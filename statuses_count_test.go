package weibo

import (
	"os"
	"testing"
)

func TestStatusesCount(t *testing.T) {
	appkey := os.Getenv("weibo_app_key")
	appsecret := os.Getenv("weibo_app_secret")
	username := os.Getenv("weibo_username")
	passwd := os.Getenv("weibo_passwd")
	redirecturi := os.Getenv("weibo_redirect_uri")
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
	resp, err := weibo.StatusesCount(token.AccessToken, 4496926160810304)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if (*resp)[0].Comments < 1 {
		t.Error("resp Comments count error")
	}
}
