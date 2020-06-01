// 登录接口

package weibo

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// CrackPinFunc 验证码破解方法类型声明
// 验证码图片以 io.Reader 类型传入，返回破解结果字符串
type CrackPinFunc func(io.Reader) (string, error)

// respPreLogin PC 端 prelogin 的返回结果
type respPreLogin struct {
	Retcode    int    `json:"retcode"`
	Servertime int    `json:"servertime"`
	Pcid       string `json:"pcid"`
	Nonce      string `json:"nonce"`
	Pubkey     string `json:"pubkey"`
	Rsakv      string `json:"rsakv"`
	IsOpenlock int    `json:"is_openlock"`
	Showpin    int    `json:"showpin"`
	Exectime   int    `json:"exectime"`
}

// respSsoLogin PC 端 ssologin 登录的返回结果
type respSsoLogin struct {
	Retcode            string   `json:"retcode"`
	Ticket             string   `json:"ticket"`
	UID                string   `json:"uid"`
	Nick               string   `json:"nick"`
	CrossDomainURLList []string `json:"crossDomainUrlList"`
}

// RespMobileLogin 移动端登录的返回结果
type RespMobileLogin struct {
	Retcode int                    `json:"retcode"`
	Msg     string                 `json:"msg"`
	Data    map[string]interface{} `json:"data"`
}

// MobileLogin 模拟移动端登录微博
// （该登录无法通过调用 Authorize 方法获取开放平台的 token）
func (w *Weibo) MobileLogin() error {
	data := url.Values{
		"username":     {w.username},
		"password":     {w.passwd},
		"savestate":    {"1"},
		"r":            {"https://weibo.cn/"},
		"ec":           {"0"},
		"pagerefer":    {"https://weibo.cn/pub/"},
		"entry":        {"mweibo"},
		"wentry":       {""},
		"loginfrom":    {""},
		"client_id":    {""},
		"code":         {""},
		"qq":           {""},
		"mainpageflag": {"1"},
		"hff":          {""},
		"hfp":          {""},
	}
	logingURL := "https://passport.weibo.cn/sso/login"
	req, err := http.NewRequest(http.MethodPost, logingURL, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin NewRequest error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "deflate, br") // no gzip
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Origin", "https://passport.weibo.cn")
	req.Header.Set("Referer", "https://passport.weibo.cn/signin/login?entry=mweibo&r=https%3A%2F%2Fweibo.cn%2F&backTitle=%CE%A2%B2%A9&vt=")
	resp, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin Do error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo MobileLogin ReadAll error")
	}
	loginResp := &RespMobileLogin{}
	if err := json.Unmarshal(body, loginResp); err != nil {
		return errors.Wrap(err, "weibo MobileLogin Unmarshal error:"+string(body))
	}
	if loginResp.Retcode != 20000000 {
		return errors.New("weibo MobileLogin loginResp Retcode error:" + string(body))
	}
	return nil
}

// RegisterCrackPinFunc 注册验证码破解方法到 Weibo 实例
// 触发验证码时自动调用注册的方法进行破解后模拟登录
func (w *Weibo) RegisterCrackPinFunc(f ...CrackPinFunc) {
	w.crackPinFuncs = append(w.crackPinFuncs, f...)
}

// 请求prelogin 获得 servertime, nonce, pubkey, rsakv 用于ssologin
func (w *Weibo) preLogin() (*respPreLogin, error) {
	// 对账号进行base64编码 对应javascript中encodeURIComponent然后base64编码
	preloginURL := "https://login.sina.com.cn/sso/prelogin.php?"
	req, err := http.NewRequest(http.MethodGet, preloginURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo PCLogin NewRequest prelogin error")
	}

	su := base64.StdEncoding.EncodeToString([]byte(w.username))
	params := url.Values{
		"entry":    {"weibo"},
		"su":       {su},
		"rsakt":    {"mod"},
		"checkpin": {"1"},
		"client":   {"ssologin.js(v1.4.19)"},
		"_":        {strconv.FormatInt(time.Now().UnixNano()/1e6, 10)},
	}
	req.URL.RawQuery = params.Encode()
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo preLogin Do prelogin error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo preLogin ReadAll prelogin error")
	}

	r := &respPreLogin{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo preLogin Unmarshal respPreLogin error:"+string(body))
	}

	if r.Retcode != 0 {
		return nil, errors.New("weibo preLogin respPreLogin Retcode error:" + string(body))
	}

	return r, nil
}

// ssologin 真正的登录
func (w *Weibo) ssoLogin(pr *respPreLogin, pinCode string) (*respSsoLogin, error) {
	ssologinURL := "https://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.19)"

	// 拼接明文js加密文件中得到
	encMsg := []byte(fmt.Sprint(pr.Servertime, "\t", pr.Nonce, "\n", w.passwd))
	// 创建公钥
	n, _ := new(big.Int).SetString(pr.Pubkey, 16)
	e, _ := new(big.Int).SetString("10001", 16)
	pubkey := &rsa.PublicKey{N: n, E: int(e.Int64())}
	// 加密公钥
	sp, err := rsa.EncryptPKCS1v15(rand.Reader, pubkey, encMsg)
	if err != nil {
		return nil, errors.Wrap(err, "weibo ssoLogin EncryptPKCS1v15 error")
	}
	// 将加密信息转换为16进制
	hexsp := hex.EncodeToString([]byte(sp))
	su := base64.StdEncoding.EncodeToString([]byte(w.username))
	data := url.Values{
		"entry":      {"account"},
		"gateway":    {"1"},
		"from":       {""},
		"savestate":  {"30"},
		"useticket":  {"1"},
		"pagerefer":  {""},
		"vsnf":       {"1"},
		"su":         {su},
		"service":    {"account"},
		"servertime": {fmt.Sprint(pr.Servertime, RandInt(1, 20))},
		"nonce":      {pr.Nonce},
		"pwencode":   {"rsa2"},
		"rsakv":      {pr.Rsakv},
		"sp":         {hexsp},
		"sr":         {"1536 * 864"},
		"encoding":   {"UTF - 8"},
		"cdult":      {"3"},
		"domain":     {"sina.com.cn"},
		"prelt":      {"95"},
		"returntype": {"TEXT"},
	}
	// 添加验证码参数
	if pinCode != "" {
		data.Add("door", pinCode)
	}
	req, err := http.NewRequest(http.MethodPost, ssologinURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "weibo PCLogin NewRequest ssologin error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo PCLogin Do ssologin error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo PCLogin ReadAll ssologin error")
	}

	// 登录结果返回结构体
	r := &respSsoLogin{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo PCLogin Unmarshal respSsoLogin error:"+string(body))
	}

	return r, nil
}

// PCLogin 模拟电脑 web 登录
// 登录后才能成功获取开放平台授权码和 token
func (w *Weibo) PCLogin() error {

	// 是否触发验证码
	hasPinCode := false
	// 验证码字符串
	pinCode := ""

	// 登录返回结果
	var sr *respSsoLogin

LOGIN: // 登录label，正常登录时跳出破解验证码的循环
	for {
		// 触发验证码时进行破解，最终得到字符串验证码
		if hasPinCode {
			// 获取验证码图片
			randNum := RandInt(10000000, 100000000) // 8位的随机数
			pinURL := fmt.Sprintf("https://login.sina.com.cn/cgi/pin.php?r=%d&s=0", randNum)
			resp, err := w.client.Get(pinURL)
			if err != nil {
				// 获取验证码图片失败则登录失败
				return errors.Wrap(err, "weibo PCLogin Get pinURL error")
			}
			defer resp.Body.Close()
			pinPic := resp.Body

			if len(w.crackPinFuncs) > 0 {
				// 有破解方法则尝试破解
				for _, crack := range w.crackPinFuncs {
					pinCode, err = crack(pinPic)
					// 破解失败尝试使用下一个crack
					if err != nil {
						logging.Error(nil, "weibo PCLogin crack pin error", zap.Error(err))
						continue
					}
					// 破解成功推出破解循环
					if pinCode != "" {
						break
					}
				}
			}
			// 破解失败，人工处理
			if err != nil || pinCode == "" {
				fmt.Println("无法自动处理登录验证码，需人工处理. ")
				// 没有注册破解方法时，保存图片,用户人工识别
				pinFilename := path.Join(os.TempDir(), "weibo_pin.png")
				pinFile, err := os.Create(pinFilename)
				if err != nil {
					return errors.Wrap(err, "weibo PCLogin Create png file error")
				}
				defer pinFile.Close()
				defer os.Remove(pinFilename)
				if _, err := io.Copy(pinFile, pinPic); err != nil {
					return errors.Wrap(err, "weibo PCLogin Copy pinPic to pinFile error")
				}
				// 尝试直接打开图片
				if err := TerminalOpen(pinFilename); err != nil {
					logging.Error(nil, "weibo TerminalOpen error", zap.Error(err))
				}

				// 等待用户输入验证码
				fmt.Printf("请输入 %s 中的验证码:", pinFilename)
				if _, err := fmt.Scanln(&pinCode); err != nil {
					return errors.Wrap(err, "weibo Scanln pinCode error")
				}
				fmt.Println("正在登录...")
			}
		}

		// 请求prelogin
		pr, err := w.preLogin()
		if err != nil {
			return err
		}

		// 请求ssologin
		sr, err = w.ssoLogin(pr, pinCode)
		if err != nil {
			return err
		}

		switch sr.Retcode {
		case "0":
			// 成功登录跳出循环
			break LOGIN
		case "101":
			// 账号密码错误
			return errors.New("weibo PCLogin respSsoLogin Retcode 101, username or password error")
		case "2070":
			// 验证码错误
			return errors.New("weibo PCLogin respSsoLogin Retcode 2070, pin code error")
		case "4049":
			// 触发验证码登录
			hasPinCode = true
		default:
			// 其他错误
			return fmt.Errorf("weibo PCLogin respSsoLogin Retcode error. %+v", sr)
		}
	}
	return w.loginSucceed(sr)
}

// loginSucceed 检查PCLogin后是否成功登录微博
func (w *Weibo) loginSucceed(resp *respSsoLogin) error {
	// 请求login_url和home_url, 进一步验证登录是否成功
	s := strings.Split(strings.Split(resp.CrossDomainURLList[0], "ticket=")[1], "&ssosavestate=")
	loginURL := fmt.Sprintf("https://passport.weibo.com/wbsso/login?ticket=%s&ssosavestate=%s&callback=sinaSSOController.doCrossDomainCallBack&scriptId=ssoscript0&client=ssologin.js(v1.4.19)&_=%s", s[0], s[1], strconv.FormatInt(time.Now().UnixNano()/1e6, 10))
	req, err := http.NewRequest(http.MethodGet, loginURL, nil)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed NewRequest loginURL error")
	}
	req.Header.Set("User-Agent", w.userAgent)
	res, err := w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed Do loginURL error")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed ReadAll loginURL error")
	}
	reg := regexp.MustCompile(`"uniqueid":"(.*?)"`)
	result := reg.FindAllStringSubmatch(string(body), -1)
	if len(result) == 0 {
		return errors.New("weibo loginSucceed uniqueid pattern not match")
	}
	uid := result[0][1]
	homeURL := fmt.Sprintf("https://weibo.com/u/%s/home", uid)
	req, err = http.NewRequest(http.MethodGet, homeURL, nil)
	req.Header.Set("User-Agent", w.userAgent)
	res, err = w.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed Do homeURL error")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "weibo loginSucceed ReadAll homeURL error")
	}
	if !strings.Contains(string(body), "我的首页") {
		return errors.New("weibo loginSucceed login failed")
	}
	return nil
}
