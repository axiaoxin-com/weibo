package weibo

import (
	"testing"
)

func TestUsersDomainShow(t *testing.T) {
	resp, err := weiboT.UsersDomainShow(tokenT, "minisnakeashin")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
}
