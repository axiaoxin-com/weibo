package weibo

import (
	"testing"
)

func TestCommentsTimeline(t *testing.T) {
	sinceID := int64(0)
	maxID := int64(0)
	count := 50
	page := 1
	trimUser := 0
	resp, err := weiboT.CommentsTimeline(tokenT, sinceID, maxID, count, page, trimUser)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
