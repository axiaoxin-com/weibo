package weibo

import (
	"testing"
)

func TestUsersShow(t *testing.T) {
	// uid := int64(1739356367)
	resp, err := weiboT.UsersShow(tokenT, 0, "陈v博")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
