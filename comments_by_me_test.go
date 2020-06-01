package weibo

import (
	"testing"
)

func TestCommentsByMe(t *testing.T) {
	sinceID := int64(0)
	maxID := int64(0)
	count := 50
	page := 1
	filterBySource := 0
	resp, err := weiboT.CommentsByMe(tokenT, sinceID, maxID, count, page, filterBySource)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
