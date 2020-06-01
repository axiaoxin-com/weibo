package weibo

import (
	"testing"
)

func TestCommentsReply(t *testing.T) {
	c, err := weiboT.CommentsByMe(tokenT, 0, 0, 1, 1, 0)
	if err != nil {
		t.Fatal("CommentsByMe err:", err)
	}
	t.Logf("%+v", c)
	cid := c.Comments[0].ID
	id := c.Comments[0].Status.ID
	resp, err := weiboT.CommentsReply(tokenT, cid, id, "爱老虎油", 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if resp.ID == 0 {
		t.Error("reply comments failed")
	}
}
