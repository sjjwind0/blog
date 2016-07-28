package login

import (
	// "errors"
	"fmt"
	"framework/config"
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

const sWeiboRedirectUrl = "http%3A%2F%2Fblog.windy.live%2Flogin%3Ftype%3Dweibo"

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
	l.appKey = config.GetDefaultConfigFileManager().ReadConfig("blog.account.weibo.key").(string)
	l.appSecret = config.GetDefaultConfigFileManager().ReadConfig("blog.account.weibo.secret").(string)
}

func (l *loginByWeb) getTokenByCode(code string) (string, string, error) {
	weburl := fmt.Sprintf(sWeiboOauth2AccessTokenUrl, l.appKey, l.appSecret, sWeiboRedirectUrl, code)
	fmt.Println("weibo START")
	fmt.Println(weburl)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", weburl, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	fmt.Printf("%+v\n", req)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	fmt.Println(string(data), err)
	var accessToken string
	var uid string
	parse := config.NewConfigContentManager(string(data))
	// accessToken = js.Get("access_token").MustString()
	accessToken = parse.ReadConfig("access_token").(string)
	uid = parse.ReadConfig("uid").(string)
	fmt.Println("test")
	return accessToken, uid, err
}

func (l *loginByWeb) getUserInfo(accessToken string, openid string) (*info.UserInfo, error) {
	fmt.Println("getUserInfo")
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
	fmt.Println(string(body))
	parse := config.NewConfigContentManager(string(body))
	var info info.UserInfo
	info.UserName = parse.ReadConfig("name").(string)
	info.Sex = parse.ReadConfig("gender").(string)
	info.UserOpenID = openid
	info.SmallFigureurl = parse.ReadConfig("profile_image_url").(string)
	return &info, nil
}

func (l *loginByWeb) Login(code string) (*info.UserInfo, error) {
	accessToken, openid, err := l.getTokenByCode(code)
	fmt.Println("get tokon start")
	if err != nil {
		return nil, err
	}
	userinfo, err := l.getUserInfo(accessToken, openid)
	fmt.Println(accessToken)
	fmt.Println(openid)
	// l.getTokenByCode(code)
	if err != nil {
		return nil, err
	}
	return userinfo, nil
}
