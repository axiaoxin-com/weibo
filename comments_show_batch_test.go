package weibo

import (
	"os"
	"testing"
)

func TestCommentsShowBatch(t *testing.T) {
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
	c, err := weibo.CommentsByMe(token.AccessToken, 0, 0, 10, 1, 0)
	if err != nil {
		t.Fatal("CommentsByMe err:", err)
	}
	cids := []int64{}
	for _, i := range c.Comments {
		cids = append(cids, i.ID)
	}
	resp, err := weibo.CommentsShowBatch(token.AccessToken, cids...)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if len(*resp) != len(cids) {
		t.Error("Comments len error")
	}
}
