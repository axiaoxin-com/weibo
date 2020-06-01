package weibo

import (
	"testing"
)

func TestCommentsShowBatch(t *testing.T) {
	c, err := weiboT.CommentsByMe(tokenT, 0, 0, 10, 1, 0)
	if err != nil {
		t.Fatal("CommentsByMe err:", err)
	}
	cids := []int64{}
	for _, i := range c.Comments {
		cids = append(cids, i.ID)
	}
	resp, err := weiboT.CommentsShowBatch(tokenT, cids...)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if len(*resp) != len(cids) {
		t.Error("Comments len error")
	}
}
