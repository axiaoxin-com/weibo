package weibo

import (
	"testing"
)

func TestCommentsCreate(t *testing.T) {
	tl, err := weiboT.StatusesHomeTimeline(tokenT, 0, 0, 1, 1, 0, 0, 1)
	if err != nil {
		t.Fatal("StatusesHomeTimeline err:", err)
	}
	t.Logf("%+v", tl)
	weiboID := tl.Statuses[0].ID
	resp, err := weiboT.CommentsCreate(tokenT, "爱老虎油", weiboID, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if resp.ID == 0 {
		t.Error("create comments failed")
	}
}
