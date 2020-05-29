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
	dom.Find("#pl_feedlist_index .card-wrap").Each(func(i int, s *goquery.Selection) {
		result := SearchWeiboResult{}

		// 获取微博 mid
		mid, _ := s.Attr("mid")
		result.Mid = mid

		// 获取用户 URL
		homePage, _ := s.Find(".avator a").Attr("href")
		if homePage != "" {
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
			picURL, _ := ss.Find("img").Attr("src")
			if picURL != "" {
				// 替换缩略图链接为大图
				picURL = strings.Replace(picURL, "orj360", "large", 1)
				picURL = strings.Replace(picURL, "thumb150", "large", 1)
				picURLs = append(picURLs, "https:"+picURL)
			}
		})
		result.Status.Origin.PicURLs = picURLs

		// 获取微博视频链接：从属性参数中正则解析
		videoActionData, _ := s.Find(".content .thumbnail .WB_video_h5").Attr("action-data")
		if videoActionData != "" {
			re, _ := regexp.Compile("video_src=(.+video)")
			matched := re.FindStringSubmatch(videoActionData)
			if len(matched) == 2 {
				u, _ := url.QueryUnescape(matched[1])
				result.Status.Origin.Video.URL = "https:" + u
			}
		}

		// 获取微博视频封面链接
		videoCover, _ := s.Find("div[node-type=fl_h5_video_pre] img").Attr("src")
		if videoCover != "" {
			result.Status.Origin.Video.CoverURL = videoCover
		}

		// 获取微博发送时间和来源
		// TODO: 时间标准化 https://github.com/dataabc/weibo-search/blob/master/weibo/utils/util.py#L53
		postTime, source := parseFromDom(s.Find(".content>.from"))
		result.Status.Origin.PostTime = postTime
		result.Status.Origin.Source = source

		// 获取微博转发数
		repost := strings.TrimSpace(s.Find(".card-act ul li:nth-of-type(2) a").Text())
		repostCount := 0
		if repost != "" {
			sl := strings.Split(repost, " ")
			if len(sl) == 2 {
				repostCount, _ = strconv.Atoi(sl[1])
			}
		}
		result.Status.Origin.RepostCount = repostCount

		// 获取微博评论数
		comment := strings.TrimSpace(s.Find(".card-act ul li:nth-of-type(3) a").Text())
		commentCount := 0
		if comment != "" {
			sl := strings.Split(comment, " ")
			if len(sl) == 2 {
				commentCount, _ = strconv.Atoi(sl[1])
			}
		}
		result.Status.Origin.CommentCount = commentCount

		// 获取微博点赞数
		like := strings.TrimSpace(s.Find(".card-act ul li:nth-of-type(4) a em").Text())
		likeCount := 0
		if like != "" {
			likeCount, _ = strconv.Atoi(like)
		}
		result.Status.Origin.LikeCount = likeCount

		// 获取微博转发内容的文本内容
		forwardContent, _ := s.Find(".content .card-comment .con div[node-type=feed_list_forwardContent] .txt").Html()
		result.Status.Forward.Content = strings.TrimSpace(forwardContent)

		// 获取微博转发内容的图片链接
		forwardPicURLs := []string{}
		s.Find(".content .card-comment .con div[node-type=feed_list_media_prev] ul li").Each(func(ii int, ss *goquery.Selection) {
			picURL, _ := ss.Find("img").Attr("src")
			if picURL != "" {
				// 替换缩略图链接为大图
				picURL = strings.Replace(picURL, "orj360", "large", 1)
				picURL = strings.Replace(picURL, "thumb150", "large", 1)
				forwardPicURLs = append(forwardPicURLs, "https:"+picURL)
			}
		})
		result.Status.Forward.PicURLs = forwardPicURLs

		// 获取转发微博的发送时间
		// TODO: 时间标准化 https://github.com/dataabc/weibo-search/blob/master/weibo/utils/util.py#L53
		forwardPostTime := strings.TrimSpace(s.Find(".content .card-comment .con .func .from a:first-of-type").Text())
		result.Status.Forward.PostTime = forwardPostTime

		// 获取转发微博的发送来源
		result.Status.Forward.Source = strings.TrimSpace(s.Find(".content .card-comment .con .func .from a:last-of-type").Text())

		// 获取转发微博的转发数
		forwardRepost := strings.TrimSpace(s.Find(".content .card-comment .con .func ul li:nth-of-type(1) a").Text())
		forwardRepostCount := 0
		if forwardRepost != "" {
			sl := strings.Split(forwardRepost, " ")
			if len(sl) == 2 {
				forwardRepostCount, _ = strconv.Atoi(sl[1])
			}
		}
		result.Status.Forward.RepostCount = forwardRepostCount

		// 获取转发微博的评论数
		forwardComment := strings.TrimSpace(s.Find(".content .card-comment .con .func ul li:nth-of-type(2) a").Text())
		forwardCommentCount := 0
		if forwardComment != "" {
			sl := strings.Split(forwardComment, " ")
			if len(sl) == 2 {
				forwardCommentCount, _ = strconv.Atoi(sl[1])
			}
		}
		result.Status.Forward.CommentCount = forwardCommentCount

		// 获取转发微博的点赞数
		forwardLike := strings.TrimSpace(s.Find(".content .card-comment .con .func ul li:nth-of-type(3) a em").Text())
		forwardLikeCount := 0
		if forwardLike != "" {
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
	html, _ := s.Html()

	// 没有来源信息的情况只返回时间
	if !strings.Contains(html, "来自") {
		postTime = strings.TrimSpace(s.Find("a:first-of-type").Text())
		return
	}

	// 05月25日 21:26 @Mayday瑪莎 转发过  来自 微博 weibo.com
	if strings.Contains(html, "转发过  来自") {
		postTime = strings.TrimSpace(s.Find("a:first-of-type").Text())
		source = strings.TrimSpace(s.Find("a:last-of-type").Text())
		return
	}

	// 05月25日 21:02 @阿信 赞过
	if strings.Contains(html, "赞过") && !strings.Contains(html, "来自") {
		postTime = strings.TrimSpace(s.Find("a:first-of-type").Text())
		return
	}

	// 05月25日 22:05 @Mayday瑪莎 转赞过  来自 微博 weibo.com
	if strings.Contains(html, "赞过") && strings.Contains(html, "来自") {
		postTime = strings.TrimSpace(s.Find("a:first-of-type").Text())
		source = strings.TrimSpace(s.Find("a:last-of-type").Text())
		return
	}

	// 8分钟前 转赞人数超过100  来自 HUAWEI P20 Pro
	if strings.Contains(html, "转赞人数超过") {
		sp := strings.Split(strings.TrimSpace(s.Find("a:first-of-type").Text()), "转赞人数超过")
		postTime = strings.TrimSpace(sp[0])
		source = strings.TrimSpace(s.Find("a:last-of-type").Text())
		return
	}

	// 其余默认视为 15分钟前  来自 iPhone客户端
	postTime = strings.TrimSpace(s.Find("a:first-of-type").Text())
	source = strings.TrimSpace(s.Find("a:last-of-type").Text())
	return
}

// SearchWeiboCondition 高级搜索筛选条件
type SearchWeiboCondition struct {
	Result string // 保存最后结果
}

// String 转化为 url 参数
func (c *SearchWeiboCondition) String() string {
	return c.Result
}

// TypeAll 设置微博搜索类型为 全部
func (c *SearchWeiboCondition) TypeAll() *SearchWeiboCondition {
	c.Result += "&typeall=1"
	return c
}

// TypeHot 设置微博搜索类型为 热门
func (c *SearchWeiboCondition) TypeHot() *SearchWeiboCondition {
	c.Result += "&xsort=hot"
	return c
}

// TypeOri 设置微博搜索类型为 原创
func (c *SearchWeiboCondition) TypeOri() *SearchWeiboCondition {
	c.Result += "&scope=ori"
	return c
}

// TypeAtten 设置微博搜索类型为 关注人
func (c *SearchWeiboCondition) TypeAtten() *SearchWeiboCondition {
	c.Result += "&atten=1"
	return c
}

// TypeVip 设置微博搜索类型为 认证用户
func (c *SearchWeiboCondition) TypeVip() *SearchWeiboCondition {
	c.Result += "&vip=1"
	return c
}

// TypeCategory 设置微博搜索类型为 认证用户
func (c *SearchWeiboCondition) TypeCategory() *SearchWeiboCondition {
	c.Result += "&category=4"
	return c
}

// TypeViewpoint 设置微博搜索类型为 认证用户
func (c *SearchWeiboCondition) TypeViewpoint() *SearchWeiboCondition {
	c.Result += "&viewpoint=1"
	return c
}

// ContainAll 设置包含条件为 全部
func (c *SearchWeiboCondition) ContainAll() *SearchWeiboCondition {
	c.Result += "&suball=1"
	return c
}

// ContainPic 设置包含条件为 包含图片
func (c *SearchWeiboCondition) ContainPic() *SearchWeiboCondition {
	c.Result += "&haspic=1"
	return c
}

// ContainVideo 设置包含条件为 包含视频
func (c *SearchWeiboCondition) ContainVideo() *SearchWeiboCondition {
	c.Result += "&hasvideo=1"
	return c
}

// ContainMusic 设置包含条件为 包含音乐
func (c *SearchWeiboCondition) ContainMusic() *SearchWeiboCondition {
	c.Result += "&hasmusic=1"
	return c
}

// ContainLink 设置包含条件为 包含短链
func (c *SearchWeiboCondition) ContainLink() *SearchWeiboCondition {
	c.Result += "&haslink=1"
	return c
}

// TimeScope 设置起止时间范围
// 时间格式：2020-05-01-18 Y-m-d-H
func (c *SearchWeiboCondition) TimeScope(begin, end string) *SearchWeiboCondition {
	c.Result += ("&timescope=custom:" + begin + ":" + end)
	return c
}

// Region 设置地点范围，传入中文
func (c *SearchWeiboCondition) Region(prov, city string) *SearchWeiboCondition {
	code := GetSearchRegionCode(prov, city)
	c.Result += ("&region=custom:" + code)
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
	URL := fmt.Sprintf("https://s.weibo.com/weibo?q=%s&page=%d%s", keyword, page, condition.String())
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
