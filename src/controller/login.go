package controller

import (
	"fmt"
	"framework/config"
	"io/ioutil"
	"net/http"
	"strings"
)

const kQQOAuthURL = "https://graph.qq.com/oauth2.0/me"

type LoginController struct {
}

func NewLoginController() *LoginController {
	return &LoginController{}
}

func (i *LoginController) Path() interface{} {
	return []string{"/login"}
}

func (l *LoginController) getOpenId(accessToken string) (string, string, error) {
	url := fmt.Sprintf("%s?access_token=%s", accessToken)
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
	begin := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	jsonBody := response[begin:end]
	fmt.Println(jsonBody)
	c := config.NewConfigContentManager(jsonBody)
	clientId := c.ReadConfig("client_id").(string)
	openId := c.ReadConfig("open_id").(string)
	fmt.Println("client_id: ", clientId)
	fmt.Println("open_id: ", openId)
	fmt.Println(string(body))
	return clientId, openId, err
}

func (l *LoginController) HandlerAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	loginType := r.Form.Get("type")
	switch loginType {
	case "qq":
		fmt.Println("login from qq")
		fmt.Println("r.Form: ", r.Form)
		accessToken := r.Form.Get("access_token")
		fmt.Println("access_token: ", accessToken)
		if len(accessToken) != 0 {
			fmt.Println(l.getOpenId(accessToken))
			// https://graph.qq.com/oauth2.0/me?access_token=47D98E8D58D74E553D9A7566C35F3C1A

		}
	case "weibo":
		fmt.Println("login from weibo")
	}
	fmt.Fprintf(w, "success")
}
