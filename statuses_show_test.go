package weibo

import (
	"testing"
)

func TestStatusesShow(t *testing.T) {
	id := int64(4409293824089668)
	resp, err := weiboT.StatusesShow(tokenT, id)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
