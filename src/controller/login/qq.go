package login

import (
	"errors"
	"fmt"
	"framework/base/config"
	"framework/base/json"
	"info"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

const kQQOAuthMeURL = "https://graph.qq.com/oauth2.0/me?access_token=%s"
const kQQOAuthTokenURL = "https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&state=%s&redirect_uri=%s"
const kQQAPIGetUserInfo = "https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s"
const kQQRedirectURL = "https%3A%2F%2Fwindyx.com%2Flogin%3Ftype%3Dqq"

var qqLoginOnce sync.Once
var qqLoginInstance *loginByQQ = nil

func GetQQLoginInstance() *loginByQQ {
	qqLoginOnce.Do(func() {
		qqLoginInstance = &loginByQQ{}
		qqLoginInstance.init()
	})
	return qqLoginInstance
}

type loginByQQ struct {
	appKey    string
	appSecret string
}

func (l *loginByQQ) init() {
	l.appKey = config.GetDefaultConfigJsonReader().Get("account.open.qq.key").(string)
	l.appSecret = config.GetDefaultConfigJsonReader().Get("account.open.qq.secret").(string)
}

func (l *loginByQQ) getTokenByCode(code string) (string, string, error) {
	// 1. 访问kQQOAuthTokenURL获取token
	url := fmt.Sprintf(kQQOAuthTokenURL, l.appKey, l.appSecret, code, "login", kQQRedirectURL)
	fmt.Println("url: ", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	response := string(body)
	fmt.Println(string(body))

	// 2. 解析返回的数据
	values := strings.Split(response, "&")
	var accessToken string
	var refreshToken string
	for _, value := range values {
		keyValue := strings.Split(value, "=")
		if len(keyValue) != 2 {
			return "", "", errors.New("return value error")
		}
		switch keyValue[0] {
		case "access_token":
			accessToken = keyValue[1]
		case "refresh_token":
			refreshToken = keyValue[1]
		}
	}
	return accessToken, refreshToken, err
}

func (l *loginByQQ) getOpenId(accessToken string) (string, error) {
	// 1. 返回kQQOAuthMeURL，获取open id
	url := fmt.Sprintf(kQQOAuthMeURL, accessToken)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	response := string(body)
	fmt.Println(string(body))

	// 2. 解析response
	begin := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	jsonBody := response[begin : end+1]
	fmt.Println(jsonBody)
	c := json.NewJsonReader(jsonBody)
	openId := c.Get("openid").(string)
	return openId, nil
}

func (l *loginByQQ) getUserInfo(accessToken string, openId string) (*info.UserInfo, error) {
	// 1. 返回kQQOAuthMeURL，获取open id
	url := fmt.Sprintf(kQQAPIGetUserInfo, accessToken, l.appKey, openId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := string(body)
	fmt.Println(string(body))

	// 2. 解析
	c := json.NewJsonReader(response)
	code := c.Get("ret").(int64)
	if code != 0 {
		fmt.Println("get user info error with error code: ", code)
		return nil, errors.New("ret error")
	}
	var info info.UserInfo
	info.UserName = c.Get("nickname").(string)
	info.Sex = c.Get("gender").(string)
	info.UserOpenID = openId
	info.SmallFigureurl = c.Get("figureurl_1").(string)
	info.BigFigureurl = c.Get("figureurl_2").(string)
	fmt.Println("userName: ", info.UserName)
	return &info, nil
}

func (l *loginByQQ) Login(code string) (*info.UserInfo, error) {
	accessToken, refreshToken, err := l.getTokenByCode(code)
	if err != nil {
		fmt.Println("getTokenByCode error: ", err)
		return nil, err
	}
	fmt.Println("accessToken: ", accessToken)
	fmt.Println("refreshToken: ", refreshToken)
	openId, err := l.getOpenId(accessToken)
	if err != nil {
		fmt.Println("getOpenId error: ", err)
		return nil, err
	}
	fmt.Println("openId: ", openId)
	userInfo, err := l.getUserInfo(accessToken, openId)
	if err != nil {
		fmt.Println("getUserInfo error: ", err)
		return nil, err
	}
	return userInfo, err
}
