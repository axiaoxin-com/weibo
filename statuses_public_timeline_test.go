package weibo

import (
	"testing"
)

func TestStatusesPublicTimeline(t *testing.T) {
	count := 3
	page := 1
	baseApp := 0
	resp, err := weiboT.StatusesPublicTimeline(tokenT, count, page, baseApp)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
