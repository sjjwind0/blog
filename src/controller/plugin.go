package controller

import (
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
	IsHtml                   bool
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
		plugin.SharePluginMgrInstance().HandleRequest(pluginId, w, r)
	} else if r.URL.Path == "/plugin" {
		p.SessionController.HandlerRequest(p, w, r)
		id, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
			return
		}
		pluginInfo, err := model.SharePluginModel().FetchPluginByPluginID(id)
		if err != nil || pluginInfo == nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
			return
		}
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
		render.IsHtml = pluginInfo.PluginType == info.PluginType_H5
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
		err = t.Execute(w, render)
		if err != nil {
			fmt.Println("execute error: ", err)
		}
	}
}
