package weibo

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// QRLogin 扫码登录
func (w *Weibo) QRLogin() error {
	// 获取二维码 url
	callback := fmt.Sprint("STK_", time.Now().UnixNano()/1e5)
	qrcodeURL := fmt.Sprint("https://login.sina.com.cn/sso/qrcode/image?entry=homepage&size=128&callback=", callback)
	resp, err := w.client.Get(qrcodeURL)
	if err != nil {
		return errors.Wrap(err, "get qrcode url error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo get qrcode url ReadAll error")
	}
	reg, _ := regexp.Compile(`"image":"(.+?)"`)
	result := reg.FindStringSubmatch(string(body))
	imgurl := result[1]
	if imgurl == "" {
		return errors.New("regexp find image failed")
	}
	// remove backslash
	imgurl = strings.Replace(imgurl, "\\", "", -1)
	reg, _ = regexp.Compile(`"qrid":"(.+?)"`)
	result = reg.FindStringSubmatch(string(body))
	qrid := result[1]
	if qrid == "" {
		return errors.New("regexp find qrid failed")
	}
	// 获取二维码图片
	respimg, err := w.client.Get(imgurl)
	if err != nil {
		return errors.Wrap(err, "get qrcode image error")
	}
	defer respimg.Body.Close()
	img := respimg.Body
	imgFilename := path.Join(os.TempDir(), "weibo_login_qr.jpg")
	imgFile, err := os.Create(imgFilename)
	if err != nil {
		return errors.Wrap(err, "save qrcode image error")
	}
	defer imgFile.Close()
	defer os.Remove(imgFilename)
	if _, err := io.Copy(imgFile, img); err != nil {
		return errors.Wrap(err, "qrcode image io copy error")
	}

	TerminalOpen(imgFilename)

	// 等待扫码
	fmt.Printf("扫描二维码: %s 进行登录...\n", imgFilename)
	alt := ""
	for {
		callback = fmt.Sprint("STK_", time.Now().UnixNano()/1e5)
		checkURL := fmt.Sprintf("https://login.sina.com.cn/sso/qrcode/check?entry=weibo&qrid=%s&callback=%s", qrid, callback)
		resp, err := w.client.Get(checkURL)
		if err != nil {
			return errors.Wrap(err, "get check url error")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "weibo get check ReadAll error")
		}
		// window.STK_16060496854274 && STK_16060496854274({"retcode":20000000,"msg":"succ","data":{"alt":"ALT-MTczOTM1NjM2Nw==-1606049685-gz-F327AC15E2A09F055157CF993714CEC6-1"}});
		if strings.Contains(string(body), "succ") {
			reg, _ = regexp.Compile(`"alt":"(.+?)"`)
			result = reg.FindStringSubmatch(string(body))
			alt = result[1]
			if alt == "" {
				return errors.New("regexp find alt failed")
			}
			break
		}
		time.Sleep(800 * time.Millisecond)
	}

	// 登录
	loginURL := fmt.Sprintf("http://login.sina.com.cn/sso/login.php?entry=weibo&returntype=TEXT&crossdomain=1&cdult=3&domain=weibo.com&savestate=30&alt=%s", alt)
	resp, err = w.client.Get(loginURL)
	if err != nil {
		return errors.Wrap(err, "qr login error")
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "weibo qr login ReadAll error")
	}
	// 登录结果返回结构体
	r := &RespQRLogin{}
	if err := json.Unmarshal(body, r); err != nil {
		return errors.Wrap(err, "weibo QRLogin Unmarshal RespQRLogin error:"+string(body))
	}

	if len(r.CrossDomainURLList) == 0 {
		return errors.New("login check not return url list:" + string(body))
	}

	// skip ssl x509 verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	w.client.Transport = tr

	for _, domainURL := range r.CrossDomainURLList {
		_, err := w.client.Get(domainURL)
		if err != nil {
			return errors.Wrap(err, "domainURL return error")
		}
	}

	return nil
}

// RespQRLogin 扫码登录返回 json 结构
type RespQRLogin struct {
	Retcode            string   `json:"retcode"`
	UID                string   `json:"uid"`
	Nick               string   `json:"nick"`
	CrossDomainURLList []string `json:"crossDomainUrlList"`
}
