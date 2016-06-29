package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"framework"
	"framework/config"
	"framework/response"
	"html/template"
	"info"
	"model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type commentRender struct {
	UserID         string
	UserName       string
	Pic            string
	CommentID      string
	CommentContent string
	CommentTime    string
	Floor          int
	ChildContent   template.HTML
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
	BlogCommentCount       string
	BlogCommentPeopleCount string
}

type BlogController struct {
	blogContentMap map[string]*[]byte
}

func NewBlogController() *BlogController {
	controller := &BlogController{}
	controller.blogContentMap = make(map[string]*[]byte)
	return controller
}

func (b *BlogController) Path() interface{} {
	return []string{"/blog", "/article", "/img"}
}

func (b *BlogController) buildCommentRender(info *info.CommentInfo, childComment *string,
	floor *int) commentRender {
	var render commentRender
	render.ChildContent = template.HTML(*childComment)
	render.CommentContent = info.Content
	render.CommentTime = "2016年12月1日 03:22"
	render.CommentID = strconv.Itoa(info.CommentID)
	render.UserID = string(info.UserID)
	render.UserName = "测试"
	render.Floor = *floor
	(*floor)++
	return render
}

func (b *BlogController) buildCommentString(child *string, info *info.CommentInfo,
	step int, floor *int) string {
	var tmpl string = ""
	if step == 0 {
		tmpl = firstComment
	} else {
		tmpl = secondComment
	}
	t, err := template.New("test").Parse(tmpl)
	buf := bytes.NewBuffer(make([]byte, 0))
	strIO := bufio.NewWriter(buf)
	if err == nil {
		t.Execute(strIO, b.buildCommentRender(info, child, floor))
	} else {
		fmt.Println(err)
	}
	strIO.Flush()
	return string(buf.Bytes())
}

func (b *BlogController) buildComment(commentTree *map[int]*info.CommentInfo,
	info *info.CommentInfo, step int, floor *int) string {
	var childComment = ""
	if info.ParentCommentID == -1 {
		return b.buildCommentString(&childComment, info, step, floor)
	}
	// 首先build 子元素
	childInfo := (*commentTree)[info.ParentCommentID]
	childComment = b.buildComment(commentTree, childInfo, step+1, floor)
	return b.buildCommentString(&childComment, info, step, floor)
}

func (b *BlogController) FetchCommentContent(blogId int) (string, error) {
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
		var floor = 1
		rawComment += b.buildComment(&commentTree, &info, 0, &floor)
	}
	return rawComment, nil
}

func (b *BlogController) readFileContent(path string) *[]byte {
	if v, ok := b.blogContentMap[path]; ok {
		return v
	}
	fileInfo, err := os.Stat(path)
	if err == nil {
		content := make([]byte, fileInfo.Size())
		file, _ := os.Open(path)
		defer file.Close()
		file.Read(content)
		b.blogContentMap[path] = &content
		return &content
	}
	return nil
}

func (b *BlogController) readBlog(w http.ResponseWriter, blogId int) {
	uuid, err := model.ShareBlogModel().GetBlogUUIDByBlogID(blogId)
	// generate blog path
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
		return
	}
	blogPath := config.GetConfigFileManager("default.conf").ReadConfig("storage.blog").(string)
	blogPath = filepath.Join(blogPath, uuid, uuid+".html")
	blogContent := b.readFileContent(blogPath)
	w.Header().Set("Accept", "*/*")
	w.Header().Set("Content-Length", strconv.Itoa(len(*blogContent)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(*blogContent)
}

func (b *BlogController) readImg(w http.ResponseWriter, imgId string) {
	imgPath := config.GetConfigFileManager("default.conf").ReadConfig("storage.img").(string)
	imgPath = filepath.Join(imgPath, imgId)
	imgContent := b.readFileContent(imgPath)
	w.Header().Set("Accept", "*/*")
	w.Header().Set("Content-Length", strconv.Itoa(len(*imgContent)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(*imgContent)
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
	content, err := b.FetchCommentContent(blogId)
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
	render.BlogID = strconv.Itoa(blogInfo.BlogID)
	render.BlogSortType = blogInfo.BlogSortType
	render.BlogTitle = blogInfo.BlogTitle
	render.BlogTime = "2015年"
	render.BlogTag = strings.Join(blogInfo.BlogTagList, "||")
	commentCount, err := model.ShareCommentModel().FetchCommentCount(blogInfo.BlogID)
	render.BlogCommentCount = strconv.Itoa(commentCount)
	peopleCount, err := model.ShareCommentModel().FetchCommentPeopleCount(blogInfo.BlogID)
	render.BlogCommentPeopleCount = strconv.Itoa(peopleCount)
	render.BlogVisitCount = strconv.Itoa(blogInfo.BlogVisitCount)
	render.CommentContent = template.HTML(content)
	t.Execute(w, render)
}

func (b *BlogController) HandlerAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/article" {
		id, err := strconv.Atoi(r.Form["id"][0])
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
			return
		}
		b.readBlog(w, id)
	} else if r.URL.Path == "/img" {
		b.readImg(w, r.Form["id"][0])
	} else {
		id, err := strconv.Atoi(r.Form["id"][0])
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
			return
		}
		b.readBlogHtml(w, id)
	}
}
