package weibo

import (
	"testing"
)

func TestStatusesHomeTimeline(t *testing.T) {
	sinceID := int64(0)
	maxID := int64(0)
	count := 1
	page := 1
	baseApp := 0
	feature := 0
	trimUser := 0
	resp, err := weiboT.StatusesHomeTimeline(tokenT, sinceID, maxID, count, page, baseApp, feature, trimUser)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if resp.TotalNumber == 0 {
		t.Error("TotalNumber == 0")
	}
}
