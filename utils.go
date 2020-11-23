// 工具类函数定义

package weibo

import (
	"fmt"
	"math/rand"
	"net"
	"os/exec"
	"runtime"
	"time"
)

// UserAgents 模拟登录时随机选择其中的 User-Agent 设置请求头
var UserAgents []string = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.5 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.109 Safari/537.36",
}

// RandUA 随机返回一个 UserAgents 中的元素
func RandUA() string {
	rand.Seed(time.Now().Unix())
	return UserAgents[rand.Intn(len(UserAgents))]
}

// RandInt 产生指定数字范围内的随机数
func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

// RealIP 获取 IP 地址
func RealIP() string {
	ip := ""
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err.Error()
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip
}

// TerminalOpen 调用终端打开文件命令
func TerminalOpen(filePath string) error {
	fmt.Println("try to open file:" + filePath)
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", filePath).Run()
	case "windows":
		return exec.Command("start", filePath).Run()
	default:
		return exec.Command("open", filePath).Run()
	}
}
