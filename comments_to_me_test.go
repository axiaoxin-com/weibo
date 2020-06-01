package weibo

import (
	"testing"
)

func TestCommentsToMe(t *testing.T) {
	sinceID := int64(0)
	maxID := int64(0)
	count := 50
	page := 1
	filterByAuthor := 0
	filterBySource := 0
	resp, err := weiboT.CommentsToMe(tokenT, sinceID, maxID, count, page, filterByAuthor, filterBySource)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
