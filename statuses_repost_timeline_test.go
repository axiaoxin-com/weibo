package weibo

import (
	"testing"
)

func TestStatusesRepostTimeline(t *testing.T) {
	id := int64(4496926160810304)
	sinceID := int64(0)
	maxID := int64(0)
	count := 1
	page := 1
	filterByAuthor := 0
	resp, err := weiboT.StatusesRepostTimeline(tokenT, id, sinceID, maxID, count, page, filterByAuthor)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
}
