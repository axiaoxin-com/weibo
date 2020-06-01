package weibo

import (
	"testing"
)

func TestSummary(t *testing.T) {
	result, err := Summary("cate=realtimehot")
	if err != nil {
		t.Error(err)
	}
	if len(result) == 0 {
		t.Error("result len = 0")
	}
	result, err = Summary("cate=socialevent")
	if err != nil {
		t.Error(err)
	}
	if len(result) == 0 {
		t.Error("result len = 0")
	}
	result, err = weiboT.SummaryRealtimeHot()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("result len = 0")
	}
	result, err = weiboT.SummarySocialEvent()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("result len = 0")
	}
	result, err = weiboT.SummaryFriendsSearch()
	if err != nil {
		t.Fatal(err)
	}
}
