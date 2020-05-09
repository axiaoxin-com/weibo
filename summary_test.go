package weibo

import (
	"fmt"
	"os"
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
	appkey := os.Getenv("weibo_app_key")
	appsecret := os.Getenv("weibo_app_secret")
	username := os.Getenv("weibo_username")
	passwd := os.Getenv("weibo_passwd")
	redirecturi := os.Getenv("weibo_redirect_uri")
	weibo := New(appkey, appsecret, username, passwd, redirecturi)
	t.Log("PCLogin...")
	if err := weibo.PCLogin(); err != nil {
		t.Fatal(err)
	}

	result, err = weibo.SummaryRealtimeHot()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("result len = 0")
	}
	result, err = weibo.SummarySocialEvent()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("result len = 0")
	}
	result, err = weibo.SummaryFriendsSearch()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) == 0 {
		t.Fatal("result len = 0")
	}
	fmt.Println(result)
}
