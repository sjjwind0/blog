package controller

import (
	"container/list"
	"framework"
	"framework/config"
	"framework/response"
	"html/template"
	"info"
	"log"
	"model"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type blogElementRender struct {
	BlogID           int
	BlogUUID         string
	BlogTitle        string
	BlogDescription  string
	BlogSortType     string
	BlogAuthor       string
	BlogTime         string
	BlogVisitCount   int
	BlogPraiseCount  int
	BlogCommentCount int
}

type indexRender struct {
	Host     *hostRender
	BlogList []*blogElementRender
	Side     *sideRender
}

func buildBlogElementRender(info *info.BlogInfo) *blogElementRender {
	var uuid string = info.BlogUUID
	storageName := config.GetConfigFileManager("default.conf").ReadConfig("storage.blog").(string)
	descriptionPath := filepath.Join(storageName, uuid, "blog.info")
	description := config.GetConfigFileManager(descriptionPath).ReadConfig("descript").(string)
	var render blogElementRender
	render.BlogAuthor = config.GetDefaultConfigFileManager().ReadConfig("owner.name").(string)
	render.BlogTitle = info.BlogTitle
	render.BlogDescription = description
	render.BlogID = info.BlogID
	render.BlogUUID = info.BlogUUID
	render.BlogPraiseCount = info.BlogPraiseCount
	render.BlogTime = FormatTime(info.BlogTime)
	render.BlogSortType = info.BlogSortType
	commentCount, _ := model.ShareCommentModel().FetchCommentCount(info.BlogID)
	render.BlogCommentCount = commentCount
	render.BlogVisitCount = info.BlogVisitCount
	return &render
}

type IndexController struct {
}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (i *IndexController) Path() interface{} {
	return []string{"/index", "/sort", "/tag", "/date"}
}

func (i *IndexController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t, err := template.ParseFiles("./src/view/html/index.html")
	if err != nil {
		log.Println(err)
	}
	var blogList *list.List = nil
	allBlogList, err := model.ShareBlogModel().FetchAllBlog()
	switch r.URL.Path {
	case "/index":
		blogList = allBlogList
	case "/sort":
		sortType := r.Form.Get("type")
		blogList, err = model.ShareBlogModel().FetchAllBlogBySortType(sortType)
	case "/tag":
		blogList = list.New()
		tagType := r.Form.Get("type")
		if err == nil {
			for iter := allBlogList.Front(); iter != nil; iter = iter.Next() {
				v := iter.Value.(info.BlogInfo)
				for _, tag := range v.BlogTagList {
					if tag == tagType {
						blogList.PushBack(v)
					}
				}
			}
		}
	case "/date":
		t := r.Form.Get("time")
		if t == "" {
			response.JsonResponse(w, framework.ErrorParamError)
			return
		}
		tList := strings.Split(t, ".")
		if len(tList) != 2 {
			response.JsonResponse(w, framework.ErrorParamError)
			return
		}
		year, err1 := strconv.Atoi(tList[0])
		month, err2 := strconv.Atoi(tList[1])
		if err1 != nil || err2 != nil {
			response.JsonResponse(w, framework.ErrorParamError)
			return
		}
		beginTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Unix()
		if month == 12 {
			month = 1
			year++
		} else {
			month++
		}
		endTime := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Unix()
		blogList, err = model.ShareBlogModel().FetchAllBlogByTime(beginTime, endTime)
	}
	if err == nil {
		var topRender indexRender
		topRender.Side = buildSideRender(allBlogList)
		for iter := blogList.Front(); iter != nil; iter = iter.Next() {
			info := iter.Value.(info.BlogInfo)
			blogRender := buildBlogElementRender(&info)
			topRender.BlogList = append(topRender.BlogList, blogRender)
		}
		topRender.Host = buildHostRender()
		t.Execute(w, &topRender)
	} else {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
	}
}
