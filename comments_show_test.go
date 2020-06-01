package weibo

import (
	"testing"
)

func TestCommentsShow(t *testing.T) {
	id := int64(4496926160810304)
	sinceID := int64(0)
	maxID := int64(0)
	count := 2
	page := 1
	filterByAuthor := 0
	resp, err := weiboT.CommentsShow(tokenT, id, sinceID, maxID, count, page, filterByAuthor)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
