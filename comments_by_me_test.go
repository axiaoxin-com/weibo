package weibo

import (
	"os"
	"testing"
)

func TestCommentsByMe(t *testing.T) {
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
	sinceID := 0
	maxID := 0
	count := 50
	page := 1
	filterBySource := 0
	resp, err := weibo.CommentsByMe(token.AccessToken, sinceID, maxID, count, page, filterBySource)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
	if resp.TotalNumber == 0 {
		t.Error("TotalNumber == 0")
	}
}
