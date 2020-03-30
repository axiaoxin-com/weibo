// https://open.weibo.com/wiki/2/statuses/share
// access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
// status	true	string	用户分享到微博的文本内容，必须做URLencode，内容不超过140个汉字，文本中不能包含“#话题词#”，同时文本中必须包含至少一个第三方分享到微博的网页URL，且该URL只能是该第三方（调用方）绑定域下的URL链接，绑定域在“我的应用 － 应用信息 － 基本应用信息编辑 － 安全域名”里设置。
// pic	false	binary	用户想要分享到微博的图片，仅支持JPEG、GIF、PNG图片，上传图片大小限制为<5M。上传图片时，POST方式提交请求，需要采用multipart/form-data编码方式。
// rip	false	string	开发者上报的操作用户真实IP，形如：211.156.0.1。

package weibo

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

/*StatusesShare 第三方分享一条链接到微博

token 为获取到的access_token内容

status 为微博文字内容

pic 为附带的一张图片，传nil则只发文字
*/
func (w *Weibo) StatusesShare(token, status string, pic io.Reader) (*StatusesShareResp, error) {
	apiURL := "https://api.weibo.com/2/statuses/share.json"
	ip := RealIP()
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)
	if pic == nil {
		data := url.Values{
			"access_token": {token},
			"status":       {status},
			"rip":          {ip},
		}
		bodyBuf = bytes.NewBufferString(data.Encode())
	} else {
		// close pic if it's a file
		if f, ok := pic.(*os.File); ok {
			defer f.Close()
		}
		picWriter, err := writer.CreateFormFile("pic", "picname.png")
		if err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare CreateFormFile error")
		}
		if _, err := io.Copy(picWriter, pic); err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare io.Copy error")
		}

		if err := writer.WriteField("access_token", token); err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare WriteField access_token error")
		}
		if err := writer.WriteField("status", status); err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare WriteField status error")
		}
		if err := writer.WriteField("rip", ip); err != nil {
			return nil, errors.Wrap(err, "weibo StatusesShare WriteField rip error")
		}
		writer.Close() // must close before new request
	}
	req, err := http.NewRequest("POST", apiURL, bodyBuf)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShare NewRequest error")
	}
	if pic == nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShare Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShare ReadAll error")
	}
	sr := &StatusesShareResp{}
	if err := json.Unmarshal(body, sr); err != nil {
		return nil, errors.Wrap(err, "weibo StatusesShare Unmarshal error:"+string(body))
	}
	if sr.Idstr == "" {
		return nil, errors.New(string(body))
	}
	return sr, nil
}
