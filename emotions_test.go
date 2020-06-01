package weibo

import (
	"testing"
)

func TestEmotions(t *testing.T) {
	// face ani cartoon
	resp, err := weiboT.Emotions(tokenT, "ani", "cnname")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", resp)
	if len(*resp) == 0 {
		t.Error("no emotions return")
	}
}
