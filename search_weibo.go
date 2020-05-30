// 微博搜索
// 参数：
// 关键词: keyword
// 类型： type   选项-> 全部=typeall 热门=hot 原创=ori 关注人=atten 认证用户=vip 媒体=category_4 观点=viewpoint 网友讨论=discuss
// 包含： sub  选项-> 全部=suball 含图片=haspic 含视频=hasvideo 含音乐=hasmusic 含短链=haslink
// 时间： 开始stime 结束etime
// 地点： 省=prov 市=city

package weibo

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/axiaoxin-com/logging"
	"github.com/pkg/errors"
)

// 解析搜索结果
func parseSearchWeiboResult(dom *goquery.Document) []SearchWeiboResult {
	results := []SearchWeiboResult{}
	dom.Find("#pl_feedlist_index div[class=card-wrap][action-type=feed_list_item]").Each(func(i int, s *goquery.Selection) {
		result := SearchWeiboResult{}

		// 获取微博 mid
		mid, _ := s.Attr("mid")
		result.ID = mid

		// 获取用户 URL
		if homePage, _ := s.Find(".avator a").Attr("href"); homePage != "" {
			result.User.HomePage = "https:" + homePage
		}

		// 获取用户头像 URL
		avatorURL, _ := s.Find(".avator a img").Attr("src")
		result.User.AvatorURL = avatorURL

		// 获取用户昵称
		result.User.NickName = s.Find(".info div .name").Text()

		// 获取微博原始内容
		content, _ := s.Find(".content>.txt").Html()
		result.Status.Origin.Content = strings.TrimSpace(content)

		// 获取微博图片链接
		picURLs := []string{}
		s.Find(".content>div[node-type=feed_list_media_prev] ul li").Each(func(ii int, ss *goquery.Selection) {
			if picURL, _ := ss.Find("img").Attr("src"); picURL != "" {
				// 替换缩略图链接为大图
				picURL = strings.Replace(picURL, "orj360", "large", 1)
				picURL = strings.Replace(picURL, "thumb150", "large", 1)
				picURLs = append(picURLs, "https:"+picURL)
			}
		})
		result.Status.Origin.PicURLs = picURLs

		// 获取微博视频链接：从属性参数中正则解析
		if videoActionData, _ := s.Find(".content .thumbnail .WB_video_h5").Attr("action-data"); videoActionData != "" {
			re, _ := regexp.Compile("video_src=(.+video)")
			if matched := re.FindStringSubmatch(videoActionData); len(matched) == 2 {
				u, _ := url.QueryUnescape(matched[1])
				result.Status.Origin.Video.URL = "https:" + u
			}
		}

		// 获取微博视频封面链接
		if videoCover, _ := s.Find("div[node-type=fl_h5_video_pre] img").Attr("src"); videoCover != "" {
			result.Status.Origin.Video.CoverURL = videoCover
		}

		// 获取微博发送时间和来源
		postTime, source := parseFromDom(s.Find(".content>.from"))
		result.Status.Origin.PostTime = NormalizeTime(postTime)
		result.Status.Origin.Source = source

		// 获取微博转发数
		repostCount := 0
		if repost := strings.TrimSpace(s.Find(".card-act ul li:nth-of-type(2) a").Text()); repost != "" {
			if sl := strings.Split(repost, " "); len(sl) == 2 {
				repostCount, _ = strconv.Atoi(sl[1])
			}
		}
		result.Status.Origin.RepostCount = repostCount

		// 获取微博评论数
		commentCount := 0
		if comment := strings.TrimSpace(s.Find(".card-act ul li:nth-of-type(3) a").Text()); comment != "" {
			if sl := strings.Split(comment, " "); len(sl) == 2 {
				commentCount, _ = strconv.Atoi(sl[1])
			}
		}
		result.Status.Origin.CommentCount = commentCount

		// 获取微博点赞数
		likeCount := 0
		if like := strings.TrimSpace(s.Find(".card-act ul li:nth-of-type(4) a em").Text()); like != "" {
			likeCount, _ = strconv.Atoi(like)
		}
		result.Status.Origin.LikeCount = likeCount

		// 获取微博转发内容的文本内容
		forwardContent, _ := s.Find(".content .card-comment .con div[node-type=feed_list_forwardContent] .txt").Html()
		result.Status.Forward.Content = strings.TrimSpace(forwardContent)

		// 获取微博转发内容的图片链接
		forwardPicURLs := []string{}
		s.Find(".content .card-comment .con div[node-type=feed_list_media_prev] ul li").Each(func(ii int, ss *goquery.Selection) {
			if picURL, _ := ss.Find("img").Attr("src"); picURL != "" {
				// 替换缩略图链接为大图
				picURL = strings.Replace(picURL, "orj360", "large", 1)
				picURL = strings.Replace(picURL, "thumb150", "large", 1)
				forwardPicURLs = append(forwardPicURLs, "https:"+picURL)
			}
		})
		result.Status.Forward.PicURLs = forwardPicURLs

		// 获取转发微博的发送时间和来源
		forwardPostTime, forwardSource := parseFromDom(s.Find(".content .card-comment .con .func .from"))
		result.Status.Forward.PostTime = NormalizeTime(forwardPostTime)
		result.Status.Forward.Source = forwardSource

		// 获取转发微博的转发数
		forwardRepostCount := 0
		if forwardRepost := strings.TrimSpace(s.Find(".content .card-comment .con .func ul li:nth-of-type(1) a").Text()); forwardRepost != "" {
			if sl := strings.Split(forwardRepost, " "); len(sl) == 2 {
				forwardRepostCount, _ = strconv.Atoi(sl[1])
			}
		}
		result.Status.Forward.RepostCount = forwardRepostCount

		// 获取转发微博的评论数
		forwardCommentCount := 0
		if forwardComment := strings.TrimSpace(s.Find(".content .card-comment .con .func ul li:nth-of-type(2) a").Text()); forwardComment != "" {
			if sl := strings.Split(forwardComment, " "); len(sl) == 2 {
				forwardCommentCount, _ = strconv.Atoi(sl[1])
			}
		}
		result.Status.Forward.CommentCount = forwardCommentCount

		// 获取转发微博的点赞数
		forwardLikeCount := 0
		if forwardLike := strings.TrimSpace(s.Find(".content .card-comment .con .func ul li:nth-of-type(3) a em").Text()); forwardLike != "" {
			forwardLikeCount, _ = strconv.Atoi(forwardLike)
		}
		result.Status.Forward.LikeCount = forwardLikeCount
		logging.Debugf(nil, "--> %d: %+v", i, result)
		results = append(results, result)
	})
	return results
}

// 处理微博来源不统一的问题
func parseFromDom(s *goquery.Selection) (postTime string, source string) {
	// 05月25日 21:26 @Mayday瑪莎 转发过  来自 微博 weibo.com
	// 05月25日 21:02 @阿信 赞过
	// 05月25日 22:05 @Mayday瑪莎 转赞过  来自 微博 weibo.com
	// 8分钟前 转赞人数超过100  来自 HUAWEI P20 Pro
	// 8分钟前 转赞人数超过2000
	html, _ := s.Html()
	postTime = strings.TrimSpace(s.Find("a:first-of-type").Text())

	// 包含 来自 直接取最后一个 a 标签
	if strings.Contains(html, "来自") {
		source = strings.TrimSpace(s.Find("a:last-of-type").Text())
	}

	// 该情况时间和文字在同一个标签
	if strings.Contains(html, "转赞人数超过") {
		sp := strings.Split(postTime, "转赞人数超过")
		postTime = strings.TrimSpace(sp[0])
	}
	return
}

// SearchWeiboCondition 高级搜索筛选条件
type SearchWeiboCondition struct {
	URLParam string // 保存最终组合到一起的搜索条件对应的 URL 参数
}

// TypeAll 设置微博搜索类型为 全部
func (c *SearchWeiboCondition) TypeAll() *SearchWeiboCondition {
	c.URLParam += "&typeall=1"
	return c
}

// TypeHot 设置微博搜索类型为 热门
func (c *SearchWeiboCondition) TypeHot() *SearchWeiboCondition {
	c.URLParam += "&xsort=hot"
	return c
}

// TypeOri 设置微博搜索类型为 原创
func (c *SearchWeiboCondition) TypeOri() *SearchWeiboCondition {
	c.URLParam += "&scope=ori"
	return c
}

// TypeAtten 设置微博搜索类型为 关注人
func (c *SearchWeiboCondition) TypeAtten() *SearchWeiboCondition {
	c.URLParam += "&atten=1"
	return c
}

// TypeVip 设置微博搜索类型为 认证用户
func (c *SearchWeiboCondition) TypeVip() *SearchWeiboCondition {
	c.URLParam += "&vip=1"
	return c
}

// TypeCategory 设置微博搜索类型为 认证用户
func (c *SearchWeiboCondition) TypeCategory() *SearchWeiboCondition {
	c.URLParam += "&category=4"
	return c
}

// TypeViewpoint 设置微博搜索类型为 认证用户
func (c *SearchWeiboCondition) TypeViewpoint() *SearchWeiboCondition {
	c.URLParam += "&viewpoint=1"
	return c
}

// ContainAll 设置包含条件为 全部
func (c *SearchWeiboCondition) ContainAll() *SearchWeiboCondition {
	c.URLParam += "&suball=1"
	return c
}

// ContainPic 设置包含条件为 包含图片
func (c *SearchWeiboCondition) ContainPic() *SearchWeiboCondition {
	c.URLParam += "&haspic=1"
	return c
}

// ContainVideo 设置包含条件为 包含视频
func (c *SearchWeiboCondition) ContainVideo() *SearchWeiboCondition {
	c.URLParam += "&hasvideo=1"
	return c
}

// ContainMusic 设置包含条件为 包含音乐
func (c *SearchWeiboCondition) ContainMusic() *SearchWeiboCondition {
	c.URLParam += "&hasmusic=1"
	return c
}

// ContainLink 设置包含条件为 包含短链
func (c *SearchWeiboCondition) ContainLink() *SearchWeiboCondition {
	c.URLParam += "&haslink=1"
	return c
}

// TimeScope 设置起止时间范围
// 时间格式：2020-05-01-18 Y-m-d-H
func (c *SearchWeiboCondition) TimeScope(begin, end string) *SearchWeiboCondition {
	c.URLParam += ("&timescope=custom:" + begin + ":" + end)
	return c
}

// Region 设置地点范围，传入中文
func (c *SearchWeiboCondition) Region(prov, city string) *SearchWeiboCondition {
	provCode, cityCode := GetSearchRegionCode(prov, city)
	c.URLParam += fmt.Sprint("&region=custom:", provCode, ":", cityCode)
	return c
}

// SearchWeibo 微博综合搜索
// pkg 级别的搜索，未登录，无法使用高级搜索，搜索内容有限,只能看评论、转发、点赞的数量
func SearchWeibo(keyword string) ([]SearchWeiboResult, error) {
	URL := "https://s.weibo.com/weibo?q=" + keyword
	dom, err := goquery.NewDocument(URL)
	if err != nil {
		return nil, errors.Wrap(err, "Search NewDocument error")
	}
	return parseSearchWeiboResult(dom), nil
}

// SearchWeibo 微博综合搜索（登录状态）
// 支持分页，翻页时不要太快，否则会跳转安全验证页面
// 支持高级搜索
func (w *Weibo) SearchWeibo(keyword string, page int, condition *SearchWeiboCondition) ([]SearchWeiboResult, error) {
	URL := fmt.Sprintf("https://s.weibo.com/weibo?q=%s&page=%d%s", keyword, page, condition.URLParam)
	logging.Debugs(nil, "weibo SearchWeibo URL:", URL)
	resp, err := w.client.Get(URL)
	if err != nil {
		return nil, errors.Wrap(err, "weibo SearchWeibo Get error")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("weibo SearchWeibo resp.Status=" + resp.Status)
	}
	dom, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, errors.Wrap(err, "weibo SearchWeibo NewDocumentFromResponse error")
	}
	return parseSearchWeiboResult(dom), nil
}
