package weibo

import (
	"testing"
)

func TestCommentsDestroyBatch(t *testing.T) {
	tl, err := weiboT.StatusesHomeTimeline(tokenT, 0, 0, 1, 1, 0, 0, 1)
	if err != nil {
		t.Fatal("StatusesHomeTimeline err:", err)
	}
	weiboID := tl.Statuses[0].ID
	cr, err := weiboT.CommentsCreate(tokenT, "爱老虎油"+string(RandInt(0, 100)), weiboID, 1)
	if err != nil {
		t.Fatal("CommentsCreate err:", err)
	}
	resp, err := weiboT.CommentsDestroyBatch(tokenT, cr.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if len(*resp) == 0 {
		t.Error("no comments deleted")
	}
}
