package weibo

import (
	"testing"
)

func TestStatusesMentions(t *testing.T) {
	sinceID := int64(0)
	maxID := int64(0)
	count := 50
	page := 1
	filterBySource := 0
	filterByAuthor := 0
	filterByType := 0
	resp, err := weiboT.StatusesMentions(tokenT, sinceID, maxID, count, page, filterBySource, filterByAuthor, filterByType)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
