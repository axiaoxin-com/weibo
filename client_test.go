package weibo

import (
	"log"
	"os"
	"testing"
)

var weiboT *Weibo
var codeT string
var tokenT string

func TestMain(m *testing.M) {
	appkey := os.Getenv("weibo_app_key")
	appsecret := os.Getenv("weibo_app_secret")
	username := os.Getenv("weibo_username")
	passwd := os.Getenv("weibo_passwd")
	redirecturi := os.Getenv("weibo_redirect_uri")
	weiboT = New(appkey, appsecret, username, passwd, redirecturi)
	if err := weiboT.QRLogin(); err != nil {
		log.Fatal(err)
	}

	code, err := weiboT.Authorize()
	if err != nil {
		log.Fatal(err)
	}
	token, err := weiboT.AccessToken(code)
	if err != nil {
		log.Fatal(err)
	}
	tokenT = token.AccessToken

	m.Run()
}
