// https://open.weibo.com/wiki/2/emotions
// 请求参数
//   access_token	true	string	采用OAuth授权方式为必填参数，OAuth授权后获得。
//   type	false	string	表情类别，face：普通表情、ani：魔法表情、cartoon：动漫表情，默认为face。
//   language	false	string	语言类别，cnname：简体、twname：繁体，默认为cnname。

package weibo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// Emotions 获取微博官方表情的详细信息
// emotionType 表情类别，face：普通表情、ani：魔法表情、cartoon：动漫表情，默认为face
// language 语言类别，cnname：简体、twname：繁体，默认为cnname
func (w *Weibo) Emotions(token, emotionType, language string) (*EmotionsResp, error) {
	apiURL := "https://api.weibo.com/2/emotions.json"
	data := url.Values{
		"access_token": {token},
		"type":         {emotionType},
		"language":     {language},
	}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "weibo Emotions NewRequest error")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = data.Encode()
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "weibo Emotions Do error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "weibo Emotions ReadAll error")
	}
	e := &ErrorResp{}
	json.Unmarshal(body, e)
	if e.Error != "" && e.ErrorCode != 0 {
		return nil, errors.New("weibo Emotions resp error:" + e.Error)
	}
	r := &EmotionsResp{}
	if err := json.Unmarshal(body, r); err != nil {
		return nil, errors.Wrap(err, "weibo Emotions Unmarshal error:"+string(body))
	}
	return r, nil
}
