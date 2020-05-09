// 登录接口

package weibo

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

// 从 dom 中获取结果
func findResult(dom *goquery.Document) []SummaryResp {
	result := []SummaryResp{}
	dom.Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		rank := s.Find(".td-01").Text()
		keyword := s.Find(".td-02 a").Text()
		heat := s.Find(".td-02 span").Text()
		tag := s.Find(".td-03 i").Text()
		url, _ := s.Find(".td-02 a").Attr("href")
		if url != "" {
			url = "https://s.weibo.com" + url
		}
		result = append(result, SummaryResp{Rank: rank, Keyword: keyword, Heat: heat, Tag: tag, URL: url})
	})
	return result
}

// Summary 微博搜索
// pkg 级别的搜索，未登录，无法获取好友搜列表，热搜榜和要闻榜正常
func Summary(param string) ([]SummaryResp, error) {
	URL := "https://s.weibo.com/top/summary/summary?" + param
	dom, err := goquery.NewDocument(URL)
	if err != nil {
		return nil, errors.Wrap(err, "Summary NewDocument error")
	}
	return findResult(dom), nil
}

// Summary 微博搜索 for client
// param:
//   cate=realtimehot 热搜榜
//   cate=socialevent 要闻榜
//   cate=total&key=friends 好友搜
func (w *Weibo) Summary(param string) ([]SummaryResp, error) {
	URL := "https://s.weibo.com/top/summary/summary?" + param
	resp, err := w.client.Get(URL)
	if err != nil {
		return nil, errors.Wrap(err, "weibo Summary Get error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("weibo Summary resp.Status=" + resp.Status)
	}
	dom, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, errors.Wrap(err, "weibo Summary NewDocumentFromResponse error")
	}
	return findResult(dom), nil
}

// SummaryRealtimeHot 微博热搜榜
func (w *Weibo) SummaryRealtimeHot() ([]SummaryResp, error) {
	return w.Summary("cate=realtimehot")
}

// SummarySocialEvent 微博要闻榜
func (w *Weibo) SummarySocialEvent() ([]SummaryResp, error) {
	return w.Summary("cate=socialevent")
}

// SummaryFriendsSearch 微博好友搜
func (w *Weibo) SummaryFriendsSearch() ([]SummaryResp, error) {
	return w.Summary("cate=total&key=friends")
}
