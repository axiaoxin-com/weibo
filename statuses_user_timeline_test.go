package weibo

import (
	"testing"
)

func TestStatusesUserTimeline(t *testing.T) {
	uid := int64(0)
	screenName := "v-bot"
	sinceID := int64(0)
	maxID := int64(0)
	count := 1
	page := 1
	baseApp := 0
	feature := 0
	trimUser := 0
	resp, err := weiboT.StatusesUserTimeline(tokenT, uid, screenName, sinceID, maxID, count, page, baseApp, feature, trimUser)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if resp.TotalNumber == 0 {
		t.Error("TotalNumber == 0")
	}
}
