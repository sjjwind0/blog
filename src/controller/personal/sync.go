package personal

import (
	"encoding/json"
	"fmt"
	"framework"
	"framework/config"
	"framework/response"
	"framework/util/archive"
	"info"
	"io/ioutil"
	"model"
	"net/http"
	"path/filepath"
	"strings"
)

type BlogMap struct {
	Name     string
	BlogList []*info.BlogInfo
}

type SyncController struct {
}

func NewSyncController() *SyncController {
	return &SyncController{}
}

func (s *SyncController) Path() interface{} {
	return "/personal/sync"
}

func (s *SyncController) listAllBlog(w http.ResponseWriter) {
	blogList, err := model.ShareBlogModel().FetchAllBlog()
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
	}
	var blogMap map[string][]*info.BlogInfo = make(map[string][]*info.BlogInfo)
	for iter := blogList.Front(); iter != nil; iter = iter.Next() {
		info := iter.Value.(info.BlogInfo)
		blogMap[info.BlogSortType] = append(blogMap[info.BlogSortType], &info)
	}
	var render []BlogMap
	for k, v := range blogMap {
		var item BlogMap
		item.Name = k
		item.BlogList = v
		render = append(render, item)
	}
	b, err := json.Marshal(render)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
		return
	}
	fmt.Fprintf(w, string(b))
}

func (s *SyncController) uploadBlog(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	fileContent := r.MultipartForm.Value["file"][0]
	uuid := r.MultipartForm.Value["uuid"][0]
	title := r.MultipartForm.Value["title"][0]
	sort := r.MultipartForm.Value["sort"][0]
	tag := r.MultipartForm.Value["tag"][0]
	tagList := strings.Split(tag, "||")
	imgContent := r.MultipartForm.Value["img"][0]

	blogStorageFilePath := config.GetConfigFileManager("default.conf").ReadConfig("storage.blog").(string)
	imgStorageFilePath := config.GetConfigFileManager("default.conf").ReadConfig("storage.img").(string)

	blogStorageFilePath = filepath.Join(blogStorageFilePath, uuid)
	fmt.Println("blogStorageFilePath: ", blogStorageFilePath)
	// unarchive
	archive.ArchiveBufferUnderPath(fileContent, blogStorageFilePath)

	archive.ArchiveBufferToPath(imgContent, imgStorageFilePath)
	// insert blog
	isExist, err := model.ShareBlogModel().BlogIsExist(uuid)
	if err == nil {
		if !isExist {
			if model.ShareBlogModel().InsertBlog(uuid, title, sort, tagList) == nil {
				response.JsonResponse(w, framework.ErrorOK)
				return
			}
		} else {
			response.JsonResponse(w, framework.ErrorBlogExist)
			return
		}
	}
	response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
}

func (s *SyncController) HandlerAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.JsonResponse(w, framework.ErrorMethodError)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if strings.Index(contentType, "application/json") != -1 {
		// post json
		result, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			response.JsonResponse(w, framework.ErrorParamError)
			return
		}
		var f interface{}
		json.Unmarshal(result, &f)
		switch f.(type) {
		case map[string]interface{}:
			info := f.(map[string]interface{})
			if api, ok := info["type"]; ok {
				switch api.(type) {
				case string:
					switch api.(string) {
					case "list":
						s.listAllBlog(w)
						return
					}
				}
			}
		}
	} else if strings.Index(contentType, "multipart/form-data") != -1 {
		// port form data
		fmt.Println("upload Blog")
		s.uploadBlog(w, r)
	}
	response.JsonResponse(w, framework.ErrorParamError)
}
