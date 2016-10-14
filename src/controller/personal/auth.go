package personal

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"framework"
	"framework/base/config"
	"framework/response"
	"framework/server"
	"io/ioutil"
	"net/http"
)

type PersonalAuthController struct {
	server.SessionController
}

func NewPersonalAuthController() *PersonalAuthController {
	return &PersonalAuthController{}
}

func (p *PersonalAuthController) Path() interface{} {
	return "/personal/auth"
}

func (p *PersonalAuthController) SessionPath() string {
	return "/"
}

func (p *PersonalAuthController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.JsonResponse(w, framework.ErrorMethodError)
		return
	}
	p.SessionController.HandlerRequest(p, w, r)

	if status, err := p.WebSession.Get("status"); err == nil && status == "auth" {
		response.JsonResponse(w, framework.ErrorOK)
		return
	}

	result, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
		return
	}
	var f interface{}
	json.Unmarshal(result, &f)
	if m, ok := f.(map[string]interface{}); ok {
		var userName, password string
		if userName, ok = m["username"].(string); !ok {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "no username")
			return
		}
		if password, ok = m["password"].(string); !ok {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "no password")
			return
		}
		defaultConfig := config.GetDefaultConfigJsonReader()
		defaultUserName := defaultConfig.Get("account.owner.authUserName").(string)
		defaultPassword := defaultConfig.Get("account.owner.authPassword").(string)

		sign := func(password string) string {
			md5Ctx := md5.New()
			md5Ctx.Write([]byte(password))
			cipherStr := md5Ctx.Sum(nil)
			return hex.EncodeToString(cipherStr)
		}

		if userName == defaultUserName && password == sign(defaultUserName+defaultPassword) {
			p.WebSession.Set("status", "auth")
			p.ResetSessionDuration()
			response.JsonResponse(w, framework.ErrorOK)
		} else {
			response.JsonResponse(w, framework.ErrorAccountAuthError)
		}
	}
}
