package controller

import (
	"framework"
	"framework/config"
	"framework/response"
	"html/template"
	"info"
	"log"
	"model"
	"net/http"
	"path/filepath"
	"sort"
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

type tagRender struct {
	Tag   string
	Count int
}

type timeRender struct {
	Date  string
	Year  int
	Month int
}

type indexRender struct {
	BlogList     []blogElementRender
	BlogTagList  []*tagRender
	BlogTimeList []*timeRender
}

type IndexController struct {
}

func NewIndexController() *IndexController {
	return &IndexController{}
}

func (i *IndexController) Path() interface{} {
	return []string{"/index", "/"}
}

func (i *IndexController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./src/view/html/index.html")
	if err != nil {
		log.Println(err)
	}
	blogList, err := model.ShareBlogModel().FetchAllBlog()
	if err == nil {
		var topRender indexRender
		var tagMap map[string]int = make(map[string]int)
		var timeMap map[string]int64 = make(map[string]int64)
		for iter := blogList.Front(); iter != nil; iter = iter.Next() {
			info := iter.Value.(info.BlogInfo)
			for tag := range info.BlogTagList {
				tagMap[info.BlogTagList[tag]]++
			}
			var tagList []*tagRender = nil
			for k, v := range tagMap {
				tagList = append(tagList, &tagRender{k, v})
			}
			timeMap[time.Unix(info.BlogTime, 0).Format("2006年01月")] = info.BlogTime
			var blogTimeList BlogTimeList = nil
			for k, v := range timeMap {
				blogTimeList = append(blogTimeList, &BlogTime{k, v})
			}
			sort.Sort(blogTimeList)
			var blogTimeStringList []*timeRender = nil
			for i := range blogTimeList {
				renderTime := time.Unix(blogTimeList[i].Time, 0)
				blogTimeStringList = append(blogTimeStringList, &timeRender{blogTimeList[i].Tag,
					renderTime.Year(), int(renderTime.Month())})
			}
			var uuid string = info.BlogUUID
			storageName := config.GetConfigFileManager("default.conf").ReadConfig("storage.blog").(string)
			descriptionPath := filepath.Join(storageName, uuid, "blog.info")
			description := config.GetConfigFileManager(descriptionPath).ReadConfig("descript").(string)
			var render blogElementRender
			render.BlogAuthor = "风"
			render.BlogTitle = info.BlogTitle
			render.BlogDescription = description
			render.BlogID = info.BlogID
			render.BlogUUID = info.BlogUUID
			render.BlogPraiseCount = 7
			render.BlogTime = FormatTime(info.BlogTime)
			render.BlogVisitCount = 9
			render.BlogCommentCount = 10
			render.BlogSortType = info.BlogSortType
			topRender.BlogList = append(topRender.BlogList, render)
			topRender.BlogTagList = tagList
			topRender.BlogTimeList = blogTimeStringList
		}
		t.Execute(w, &topRender)
	} else {
		response.JsonResponse(w, framework.ErrorSQLError)
	}
}
