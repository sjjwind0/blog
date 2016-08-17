package controller

import (
	"fmt"
	"framework"
	"framework/config"
	"framework/response"
	"framework/server"
	"html/template"
	"info"
	"model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type blogCommentRender struct {
	UserID         string
	UserName       string
	Pic            string
	CommentID      string
	CommentContent string
	CommentTime    string
	Floor          int
	User           *info.UserInfo
	ChildContent   template.HTML
}

type userRender struct {
	IsLogin  bool
	UserID   string
	NickName string
	Pic      string
}

type blogRender struct {
	CommentContent         template.HTML
	BlogID                 string
	BlogTitle              string
	BlogTag                string
	BlogSortType           string
	BlogTime               string
	Author                 string
	BlogVisitCount         string
	BlogContent            template.HTML
	BlogCommentCount       string
	BlogCommentPeopleCount string
	User                   userRender
	Side                   *sideRender
	Host                   *hostRender
}

type BlogController struct {
	server.SessionController
	blogContentMap map[string]*[]byte
}

func NewBlogController() *BlogController {
	controller := &BlogController{}
	controller.blogContentMap = make(map[string]*[]byte)
	return controller
}

func (b *BlogController) Path() interface{} {
	return "/blog"
}

func (b *BlogController) SessionPath() string {
	return "/"
}

func (b *BlogController) fetchCommentContent(blogId int) (string, error) {
	commentList, err := model.ShareCommentModel().FetchAllCommentByBlogId(blogId)
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

func (b *BlogController) readBlogContent(blogId int) string {
	uuid, err := model.ShareBlogModel().GetBlogUUIDByBlogID(blogId)
	// generate blog path
	if err != nil {
		return ""
	}
	blogPath := config.GetDefaultConfigFileManager().ReadConfig("blog.storage.file.blog").(string)
	blogPath = filepath.Join(blogPath, uuid, uuid+".html")
	fileInfo, err := os.Stat(blogPath)
	if err == nil {
		content := make([]byte, fileInfo.Size())
		file, _ := os.Open(blogPath)
		defer file.Close()
		file.Read(content)
		return string(content)
	}
	return ""
}

func (b *BlogController) readBlogHtml(w http.ResponseWriter, blogId int) {
	if err := model.ShareBlogModel().AddVisitCount(blogId); err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorRenderError, err.Error())
		return
	}
	t, err := template.ParseFiles("./src/view/html/blog.html")
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorRenderError, err.Error())
		return
	}
	content, err := b.fetchCommentContent(blogId)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
		return
	}
	blogInfo, err := model.ShareBlogModel().FetchBlogByBlogID(blogId)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
		return
	}
	var render blogRender
	render.Host = buildHostRender()
	render.BlogID = strconv.Itoa(blogInfo.BlogID)
	render.BlogSortType = blogInfo.BlogSortType
	render.BlogTitle = blogInfo.BlogTitle
	render.BlogTime = FormatRealTime(blogInfo.BlogTime)
	render.BlogTag = strings.Join(blogInfo.BlogTagList, "||")
	commentCount, err := model.ShareCommentModel().FetchCommentCount(blogInfo.BlogID)
	render.BlogCommentCount = strconv.Itoa(commentCount)
	peopleCount, err := model.ShareCommentModel().FetchCommentPeopleCount(blogInfo.BlogID)
	render.BlogCommentPeopleCount = strconv.Itoa(peopleCount)
	render.BlogVisitCount = strconv.Itoa(blogInfo.BlogVisitCount)
	render.CommentContent = template.HTML(content)
	render.BlogContent = template.HTML(b.readBlogContent(blogId))
	render.Author = config.GetDefaultConfigFileManager().ReadConfig("blog.owner.name").(string)
	v, err := b.SessionController.WebSession.Get("status")
	if err == nil {
		if v.(string) == "login" {
			render.User.IsLogin = true
			uid, err := b.SessionController.WebSession.Get("id")
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
	blogList, err := model.ShareBlogModel().FetchAllBlog()
	if err == nil {
		render.Side = buildSideRender(blogList)
	}
	t.Execute(w, render)
}

func (b *BlogController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/blog" {
		b.SessionController.HandlerRequest(b, w, r)
		id, err := strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
			return
		}
		b.readBlogHtml(w, id)
	} else {
		response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
	}
}
