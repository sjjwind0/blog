package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"framework"
	"framework/base/config"
	"framework/response"
	"framework/server"
	"html/template"
	"info"
	"model"
	"net/http"
	"net/url"
	"plugin"
	"strconv"
	"strings"
)

type pluginRender struct {
	Host *hostRender
	*info.PluginInfo
	PluginVisitCount         string
	PluginCommentPeopleCount string
	PluginCommentCount       string
	PluginCommentContent     template.HTML
	Author                   string
	DisplayTime              string
	User                     userRender
}

type PluginController struct {
	server.SessionController
}

func NewPluginController() *PluginController {
	return &PluginController{}
}

func (p *PluginController) Path() (interface{}, bool) {
	return "/plugin", true
}

func (p *PluginController) SessionPath() string {
	return "/"
}

func (p *PluginController) fetchCommentContent(blogId int) (string, error) {
	commentList, err := model.ShareCommentModel().FetchAllCommentByBlogId(info.CommentType_Blog, blogId)
	if err != nil {
		return "", err
	}
	// 组成一个tree的形式
	var commentTree map[int]*info.CommentInfo = make(map[int]*info.CommentInfo)
	for iter := commentList.Front(); iter != nil; iter = iter.Next() {
		info := iter.Value.(info.CommentInfo)
		commentTree[info.CommentID] = &info
	}
	var rawComment string = ""
	for iter := commentList.Front(); iter != nil; iter = iter.Next() {
		info := iter.Value.(info.CommentInfo)
		rawComment += buildOneCommentFromCommentTree(&commentTree, &info)
	}
	return rawComment, nil
}

func (p *PluginController) handlePluginRequest(pluginId int, w http.ResponseWriter, r *http.Request) {
	mapStringInterfaceToMapStringString := func(value map[string]interface{}) map[string]string {
		var retMap map[string]string = make(map[string]string)
		for k, v := range value {
			retMap[k] = v.(string)
		}
		return retMap
	}
	// build normal request
	var stringArrayMerge = func(header map[string][]string) map[string]string {
		var ret map[string]string = make(map[string]string)
		for key, value := range header {
			ret[key] = value[0]
		}
		return ret
	}
	var httpRequest map[string]interface{} = map[string]interface{}{
		"method": r.Method,
		"url":    r.URL.String(),
		"header": stringArrayMerge(r.Header),
		"length": r.ContentLength,
		"host":   r.Host,
		"addr":   r.RemoteAddr,
		"form":   stringArrayMerge(r.Form),
	}
	httpRequestBytes, _ := json.Marshal(httpRequest)
	var req map[string]interface{} = map[string]interface{}{
		"type":    "com.request.http",
		"request": string(httpRequestBytes),
	}
	singal := make(chan bool)
	requestBytes, _ := json.Marshal(req)
	callback := func(code int, response string) {
		if code == 0 {
			// to http response
			var js interface{} = nil
			err := json.Unmarshal([]byte(response), &js)
			if err != nil {
				fmt.Println("ServerHTTP json.Unmarshal error: ", err)
				singal <- true
				return
			}
			if v, ok := js.(map[string]interface{}); ok {
				if header, ok := v["header"]; ok {
					nativeHeader := mapStringInterfaceToMapStringString(header.(map[string]interface{}))
					for key, value := range nativeHeader {
						w.Header().Set(key, value)
					}
				}
				if length, ok := v["length"]; ok {
					w.Header().Set("Content-Length", fmt.Sprintf("%d", int64(length.(float64))))
				}
				if code, ok := v["code"]; ok {
					w.WriteHeader(int(code.(float64)))
				}
				if body, ok := v["body"]; ok {
					nativeBody := body.(string)
					data, _ := base64.StdEncoding.DecodeString(nativeBody)
					w.Write(data)
				}
			}
		} else {
			fmt.Println("callback error: ", response)
		}
		singal <- true
	}
	plugin.SharePluginMgrInstance().CallMethod(pluginId, string(requestBytes), callback)
	<-singal
}

func (p *PluginController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if strings.HasPrefix(r.URL.Path, "/plugin/") {
		values := strings.Split(r.URL.Path, "/")
		pluginId, err := strconv.Atoi(values[2])
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
			return
		}
		rawURL := r.URL.RequestURI()
		index := strings.Index(rawURL[8:], "/")
		newUrl := rawURL[8+index:]
		r.URL, _ = url.Parse(newUrl)
		p.handlePluginRequest(pluginId, w, r)
	} else if r.URL.Path == "/plugin" {
		p.SessionController.HandlerRequest(p, w, r)
		id, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
			return
		}
		pluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(id)
		if err != nil || pluginInfo == nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
			return
		}
		plugin.SharePluginMgrInstance().LoadPlugin(id)
		t, err := template.ParseFiles("./src/view/html/plugin.html")

		var render pluginRender
		render.Host = buildHostRender()
		render.PluginInfo = pluginInfo

		render.Host = buildHostRender()

		content, err := p.fetchCommentContent(id)
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
			return
		}

		commentCount, err := model.ShareCommentModel().FetchCommentCount(
			info.CommentType_Plugin, pluginInfo.PluginID)
		render.PluginCommentCount = strconv.Itoa(commentCount)
		peopleCount, err := model.ShareCommentModel().FetchCommentPeopleCount(
			info.CommentType_Plugin, pluginInfo.PluginID)
		render.PluginCommentPeopleCount = strconv.Itoa(peopleCount)
		render.PluginVisitCount = strconv.Itoa(0)
		render.PluginCommentContent = template.HTML(content)
		render.Author = config.GetDefaultConfigJsonReader().Get("account.owner.name").(string)
		render.DisplayTime = FormatRealTime(pluginInfo.PluginTime)
		v, err := p.SessionController.WebSession.Get("status")
		if err == nil {
			if v.(string) == "login" {
				render.User.IsLogin = true
				uid, err := p.SessionController.WebSession.Get("id")
				if err == nil {
					userId, err := strconv.Atoi(uid.(string))
					userInfo, err := model.ShareUserModel().GetUserInfoById(int64(userId))
					if err == nil && userInfo != nil {
						render.User.NickName = userInfo.UserName
						render.User.Pic = userInfo.SmallFigureurl
						render.User.UserID = uid.(string)
					} else {
						render.User.IsLogin = false
					}
				} else {
					fmt.Println("err: ", err)
				}
			} else {
				render.User.IsLogin = false
			}
		} else {
			render.User.IsLogin = false
		}
		t.Execute(w, render)
	}
}
