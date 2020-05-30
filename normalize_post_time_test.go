package weibo

import (
	"testing"
	"time"
)

func TestNormalizeTime(t *testing.T) {
	now := time.Now()
	s := NormalizeTime("刚刚")
	t.Log("刚刚", s)
	if now.Format(TimeLayout) != s {
		t.Error("刚刚：", s)
	}
	s = NormalizeTime("18秒前")
	t.Log("18秒前", s)
	d, _ := time.ParseDuration("-18s")
	if now.Add(d).Format(TimeLayout) != s {
		t.Error("18秒前：", s)
	}
	s = NormalizeTime("18分钟前")
	t.Log("18分钟前", s)
	d, _ = time.ParseDuration("-18m")
	if now.Add(d).Format(TimeLayout) != s {
		t.Error("18分钟前：", s)
	}
	s = NormalizeTime("1小时前")
	t.Log("1小时前", s)
	d, _ = time.ParseDuration("-1h")
	if now.Add(d).Format(TimeLayout) != s {
		t.Error("1小时前：", s)
	}
	s = NormalizeTime("今天20:20")
	t.Log("今天20:20", s)
	if now.Format("2006年01月02日")+" 20:20" != s {
		t.Error("今天20:20：", s)
	}
	s = NormalizeTime("05月31日 20:20")
	t.Log("05月31日 20:20", s)
	if now.Format("2006年")+"05月31日 20:20" != s {
		t.Error("05月31日 20:20：", s)
	}
}
