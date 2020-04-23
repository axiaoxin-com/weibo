package weibo

import (
	"os"
	"testing"
)

func TestStatusesGo(t *testing.T) {
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
	uid := int64(5023553094)
	id := int64(4496926160810304)
	resp, err := weibo.StatusesGo(token.AccessToken, uid, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("go url: %+v", resp)
}
