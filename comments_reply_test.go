package weibo

import (
	"os"
	"testing"
)

func TestCommentsReply(t *testing.T) {
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
	c, err := weibo.CommentsByMe(token.AccessToken, 0, 0, 1, 1, 0)
	if err != nil {
		t.Fatal("StatusesHomeTimeline err:", err)
	}
	t.Logf("%+v", c)
	cid := c.Comments[0].ID
	id := c.Comments[0].Status.ID
	resp, err := weibo.CommentsReply(token.AccessToken, cid, id, "爱老虎油", 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if resp.ID == 0 {
		t.Error("reply comments failed")
	}
}
