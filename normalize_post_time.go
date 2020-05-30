// 微博数据相关标准化

package weibo

import (
	"strings"
	"time"
)

// TimeLayout 微博展示的时间格式
const TimeLayout = "2006年01月02日 15:04"

// NormalizeTime 标准化微博时间
// 刚刚、x秒前、x分钟前、x小时前、今天H:M、x月x日 -> Y-m-d H:M
func NormalizeTime(t string) (nt string) {
	now := time.Now()
	if t == "" {
		return
	} else if strings.Contains(t, "刚刚") {
		nt = now.Format(TimeLayout)
	} else if strings.Contains(t, "秒") {
		x := "-" + strings.TrimSpace(strings.Split(t, "秒")[0]) + "s"
		d, _ := time.ParseDuration(x)
		nt = now.Add(d).Format(TimeLayout)
	} else if strings.Contains(t, "分钟") {
		x := "-" + strings.TrimSpace(strings.Split(t, "分钟")[0]) + "m"
		d, _ := time.ParseDuration(x)
		nt = now.Add(d).Format(TimeLayout)
	} else if strings.Contains(t, "小时") {
		x := "-" + strings.TrimSpace(strings.Split(t, "小时")[0]) + "h"
		d, _ := time.ParseDuration(x)
		nt = now.Add(d).Format(TimeLayout)
	} else if strings.Contains(t, "今天") {
		nt = now.Format("2006年01月02日 ") + strings.TrimSpace(strings.Split(t, "今天")[1])
	} else if !strings.Contains(t, "年") {
		nt = now.Format("2006年") + strings.TrimSpace(t)
	} else {
		nt = t
	}
	return
}
