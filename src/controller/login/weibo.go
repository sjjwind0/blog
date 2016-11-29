package login

import (
	// "errors"
	"fmt"
	"framework/base/config"
	"framework/base/json"
	"info"
	"io/ioutil"
	"net/http"
	"sync"
)

type loginByWeb struct {
	appKey    string
	appSecret string
}

const sWeiboOauth2AccessTokenUrl = "https://api.weibo.com/oauth2/access_token?client_id=%s&client_secret=%s&grant_type=authorization_code&redirect_uri=%s&code=%s"

const sWeiboRedirectUrl = "https%3A%2F%2Fwindyx.com%2Flogin%3Ftype%3Dweibo"

const sWeiboUserShowUrl = "https://api.weibo.com/2/users/show.json?access_token=%s&uid=%s"

const sWeiborevokeoauth2 = "https://api.weibo.com/oauth2/revokeoauth2?access_token=%s"

var webLoginOnce sync.Once
var webLoginInstance *loginByWeb = nil

func GetWebLoginInstance() *loginByWeb {
	webLoginOnce.Do(func() {
		webLoginInstance = &loginByWeb{}
		webLoginInstance.init()
	})
	return webLoginInstance
}

func (l *loginByWeb) init() {
	l.appKey = config.GetDefaultConfigJsonReader().Get("account.open.weibo.key").(string)
	l.appSecret = config.GetDefaultConfigJsonReader().Get("account.open.weibo.secret").(string)
}

func (l *loginByWeb) getTokenByCode(code string) (string, string, error) {
	weburl := fmt.Sprintf(sWeiboOauth2AccessTokenUrl, l.appKey, l.appSecret, sWeiboRedirectUrl, code)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", weburl, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	var accessToken string
	var uid string
	parse := json.NewJsonReader(string(data))
	fmt.Println("parse: ", parse)
	accessToken = parse.Get("access_token").(string)
	uid = parse.Get("uid").(string)
	return accessToken, uid, err
}

func (l *loginByWeb) getUserInfo(accessToken string, openid string) (*info.UserInfo, error) {
	weburl := fmt.Sprintf(sWeiboUserShowUrl, accessToken, openid)
	resp, err := http.Get(weburl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	parse := json.NewJsonReader(string(body))
	var info info.UserInfo
	info.UserName = parse.Get("name").(string)
	info.Sex = parse.Get("gender").(string)
	info.UserOpenID = openid
	info.SmallFigureurl = parse.Get("profile_image_url").(string)
	return &info, nil
}

func (l *loginByWeb) Login(code string) (*info.UserInfo, error) {
	accessToken, openid, err := l.getTokenByCode(code)
	if err != nil {
		return nil, err
	}
	userinfo, err := l.getUserInfo(accessToken, openid)
	if err != nil {
		return nil, err
	}
	return userinfo, nil
}
