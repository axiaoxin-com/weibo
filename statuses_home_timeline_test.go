package weibo

import (
	"os"
	"testing"
)

func TestStatusesHomeTimeline(t *testing.T) {
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
	sinceID := int64(0)
	maxID := int64(0)
	count := 1
	page := 1
	baseApp := 0
	feature := 0
	trimUser := 0
	resp, err := weibo.StatusesHomeTimeline(token.AccessToken, sinceID, maxID, count, page, baseApp, feature, trimUser)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if resp.TotalNumber == 0 {
		t.Error("TotalNumber == 0")
	}
}
