package weibo

import (
	"os"
	"testing"
)

func TestSearchWeiboByClient(t *testing.T) {
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

	keyword := "五月天"
	page := 1
	cond := &SearchWeiboCondition{}
	// 原创 + 会员认证 + 包含音乐 + 时间段 + 地点
	// cond = cond.TypeOri().TypeVip().ContainMusic().TimeScope("2020-05-01-0", "2020-06-01-0").Region("四川", "成都")
	resp, err := weibo.SearchWeibo(keyword, page, cond)
	if err != nil {
		t.Error(err)
	}
	if len(resp) == 0 {
		t.Error("no weibo search return")
	}
}
