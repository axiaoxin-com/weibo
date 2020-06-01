package weibo

import (
	"testing"
)

func TestStatusesCount(t *testing.T) {
	resp, err := weiboT.StatusesCount(tokenT, 4496926160810304)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", resp)
	if (*resp)[0].Comments < 1 {
		t.Error("resp Comments count error")
	}
}
