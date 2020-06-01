package weibo

import (
	"testing"
)

func TestStatusesGo(t *testing.T) {
	uid := int64(5023553094)
	id := int64(4496926160810304)
	resp, err := weiboT.StatusesGo(tokenT, uid, id)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("go url: %+v", resp)
}
