package controller

import (
	"fmt"
	"framework"
	"framework/config"
	"framework/response"
	"framework/server"
	"model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ArticleController struct {
	server.SessionController
	blogContentMap map[string]*[]byte
}

func NewArticleController() *ArticleController {
	controller := &ArticleController{}
	controller.blogContentMap = make(map[string]*[]byte)
	return controller
}

func (b *ArticleController) Path() (interface{}, bool) {
	return []string{"/article", "/cover"}, true
}

func (b *ArticleController) readFileContent(path string) *[]byte {
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

func (b *ArticleController) readBlog(w http.ResponseWriter, blogId int) {
	uuid, err := model.ShareBlogModel().GetBlogUUIDByBlogID(blogId)
	// generate blog path
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
		return
	}
	blogPath := config.GetDefaultConfigFileManager().ReadConfig("blog.storage.file.blog").(string)
	blogPath = filepath.Join(blogPath, uuid, uuid+".html")
	blogContent := b.readFileContent(blogPath)
	w.Header().Set("Accept", "*/*")
	w.Header().Set("Content-Length", strconv.Itoa(len(*blogContent)))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(*blogContent)
}

func (b *ArticleController) readRes(w http.ResponseWriter, path string) {
	ext := filepath.Ext(path)
	imgContent := b.readFileContent(path)
	if imgContent == nil {
		return
	}
	w.Header().Set("Accept", "*/*")
	w.Header().Set("Content-Length", strconv.Itoa(len(*imgContent)))
	w.Header().Set("Content-Type", server.QueryContentTypeByExt(ext))
	w.Write(*imgContent)
}

func (b *ArticleController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path == "/article" {
		id, err := strconv.Atoi(r.Form["id"][0])
		if err != nil {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
			return
		}
		b.readBlog(w, id)
	} else if r.URL.Path == "/cover" {
		uuid := r.Form.Get("id")
		blogPath := config.GetDefaultConfigFileManager().ReadConfig("blog.storage.file.blog").(string)
		imgPath := filepath.Join(blogPath, uuid, "cover.jpg")
		b.readRes(w, imgPath)
	} else if strings.HasPrefix(r.URL.Path, "/article/") {
		// 首先解析uuid
		url := r.URL.Path[1:]
		part := strings.Split(url, "/")
		fmt.Println("part: ", part)
		if len(part) == 4 {
			uuid := part[1]
			blogPath := config.GetDefaultConfigFileManager().ReadConfig("blog.storage.file.blog").(string)
			resPath := filepath.Join(blogPath, uuid, "res", part[2], part[3])
			fmt.Println("resPath: ", resPath)
			b.readRes(w, resPath)
		} else {
			response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
		}
	}
}
