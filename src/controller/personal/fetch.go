package personal

import (
	"encoding/json"
	"fmt"
	"framework"
	"framework/response"
	"framework/server"
	"info"
	"io/ioutil"
	"model"
	"net/http"
)

type PersonalFetchController struct {
	server.SessionController
}

func NewPersonalFetchController() *PersonalFetchController {
	return &PersonalFetchController{}
}

func (p *PersonalFetchController) Path() interface{} {
	return "/personal/fetch"
}

func (p *PersonalFetchController) SessionPath() string {
	return "/"
}

func (p *PersonalFetchController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.JsonResponse(w, framework.ErrorMethodError)
		return
	}
	p.SessionController.HandlerRequest(p, w, r)

	if status, err := p.WebSession.Get("status"); err != nil || status != "auth" {
		fmt.Println("err: ", err.Error(), "\tstatus: ", status)
		response.JsonResponseWithMsg(w, framework.ErrorAccountAuthError, "not auth")
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
		var fetchType string
		if fetchType, ok = m["type"].(string); !ok {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "no username")
			return
		}
		switch fetchType {
		case "blog":
			blogList, err := model.ShareBlogModel().FetchAllBlog()
			if err != nil {
				response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
				return
			}
			var retBlgoList []interface{}
			for iter := blogList.Front(); iter != nil; iter = iter.Next() {
				blogInfo := iter.Value.(info.BlogInfo)
				retBlogInfo := map[string]interface{}{
					"id":   blogInfo.BlogID,
					"name": blogInfo.BlogTitle,
					"time": blogInfo.BlogTime,
					"sort": blogInfo.BlogSortType,
					"tag":  blogInfo.BlogTagList,
				}
				retBlgoList = append(retBlgoList, retBlogInfo)
			}
			response.JsonResponseWithData(w, framework.ErrorOK, "", retBlgoList)
		default:
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "unsupport type")
		}
	} else {
		response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
	}
}
