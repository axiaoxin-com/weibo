package weibo

import (
	"os"
	"testing"
)

func TestCommentsCreate(t *testing.T) {
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
	tl, err := weibo.StatusesHomeTimeline(token.AccessToken, 0, 0, 1, 1, 0, 0, 1)
	if err != nil {
		t.Fatal("StatusesHomeTimeline err:", err)
	}
	t.Logf("%+v", tl)
	weiboID := tl.Statuses[0].ID
	resp, err := weibo.CommentsCreate(token.AccessToken, "爱老虎油", weiboID, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if resp.ID == 0 {
		t.Error("create comments failed")
	}
}
