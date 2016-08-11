package personal

import (
	"fmt"
	"framework"
	"framework/config"
	"framework/response"
	"framework/server"
	"framework/util/archive"
	"io"
	"model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FileController struct {
	server.SessionController
}

func NewPersonalFileController() *FileController {
	return &FileController{}
}

func (f *FileController) Path() interface{} {
	return []string{"/personal/file-download", "/personal/file-upload"}
}

func (f *FileController) SessionPath() string {
	return "/"
}

func (f *FileController) handlerDownloadRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	blogId, err := strconv.Atoi(r.Form.Get("id"))
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
		return
	}
	// read raw zip file path
	rawPath := config.GetDefaultConfigFileManager().ReadConfig("blog.storage.file.raw").(string)
	blogInfo, err := model.ShareBlogModel().FetchBlogByBlogID(blogId)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
		return
	}
	blogPath := filepath.Join(rawPath, blogInfo.BlogUUID)
	file, err := os.Open(blogPath)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorFileNotExist, err.Error())
		return
	}
	fileInfo, err := os.Stat(blogPath)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorFileNotExist, err.Error())
		return
	}
	var content []byte = make([]byte, fileInfo.Size())
	file.Read(content)
	file.Close()
	w.Header().Set("Accept", "*/*")
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))
	w.Header().Set("Content-Disposition", "attachment; filename="+blogInfo.BlogTitle)
	w.Write(content)
}

/* 接受multi-part form格式，格式如下：
** @version 1
** 1. raw, 原始zip文件，包括所有的未经处理了的文件，服务端存储raw文件，用来供客户端下载恢复。
** 2. html, 经过处理的主要html文件。
** 3. meta信息, {"title": "xx", "tag": ["tag1", "tag2"], "sort": "xxx"}。
** 4. res, html中所需要的所有资源文件。
** 文件目录格式如下
**	raw:
**		- uuid_1.raw
**		- uuid_2.raw
**	blog:
**		- uuid
**			- uuid_1.html
**			- cover.jpg
**			- res
**				- html
**	 			- css
**	 			- js
**				- img
**	 			- font
** raw文件放在raw目录不对外开放，html文件以及res文件放在blog目录，meta信息放数据库。
 */
func (f *FileController) handlerUploadRequest(w http.ResponseWriter, r *http.Request) {
	const _24K = (1 << 20) * 24
	if err := r.ParseMultipartForm(_24K); nil != err {
		fmt.Println("r.ParseMultipartForm: ", err)
		return
	}
	checkFolder := func(path string) {
		_, err := os.Stat(path)
		if !(err == nil || os.IsExist(err)) {
			os.MkdirAll(path, 0755)
		}
	}
	saveFile := func(name string, path string) string {
		file, handler, err := r.FormFile(name)
		if err != nil {
			fmt.Println("r.FormFile: ", err)
			return ""
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile(filepath.Join(path, handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		defer f.Close()
		io.Copy(f, file)
		return handler.Filename
	}

	saveTmpPath := "/tmp"

	// 1. save raw.zip
	rawZipName := saveFile("raw", saveTmpPath)
	// 2. save web html
	webHtmlName := saveFile("web", saveTmpPath)
	// 3. save blog.info
	blogInfoName := saveFile("info", saveTmpPath)
	// 4. save res.zip
	resZipName := saveFile("res", saveTmpPath)
	// 5. save cover.jpg
	coverImgName := saveFile("img", saveTmpPath)
	// 6. read blog.info
	blogMetaInfo := config.GetConfigFileManager(filepath.Join(saveTmpPath, blogInfoName))
	uuid := blogMetaInfo.ReadConfig("uuid").(string)
	title := blogMetaInfo.ReadConfig("title").(string)
	tag := blogMetaInfo.ReadConfig("tag").(string)
	tagList := strings.Split(tag, "||")
	sort := blogMetaInfo.ReadConfig("sort").(string)
	isExist, err := model.ShareBlogModel().BlogIsExist(uuid)
	fmt.Println("uuid: ", uuid)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
		return
	}
	// 7. archive to path
	rawRootPath := config.GetDefaultConfigFileManager().ReadConfig("blog.storage.file.raw").(string)
	checkFolder(rawRootPath)
	blogRootPath := config.GetDefaultConfigFileManager().ReadConfig("blog.storage.file.blog").(string)
	blogRootPath = filepath.Join(blogRootPath, uuid)
	checkFolder(rawRootPath)

	rawZipPath := filepath.Join(rawRootPath, uuid+".zip")
	os.Rename(filepath.Join(saveTmpPath, rawZipName), rawZipPath)

	infoPath := filepath.Join(blogRootPath, blogInfoName)
	os.Rename(filepath.Join(saveTmpPath, blogInfoName), infoPath)

	webPath := filepath.Join(blogRootPath, uuid+".html")
	os.Rename(filepath.Join(saveTmpPath, webHtmlName), webPath)

	resZipPath := filepath.Join(blogRootPath, resZipName)
	os.Rename(filepath.Join(saveTmpPath, resZipName), resZipPath)

	coverImgPath := filepath.Join(blogRootPath, coverImgName)
	os.Rename(filepath.Join(saveTmpPath, coverImgName), coverImgPath)

	// archive res zip to folder
	err = archive.UnZip(resZipPath)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorParamError, err.Error())
		return
	}
	// write db
	if isExist {
		// 更新blog
		fmt.Println("update blog")
		model.ShareBlogModel().UpdateBlog(uuid, title, sort, tagList)
	} else {
		// 插入新blog
		fmt.Println("insert blog")
		model.ShareBlogModel().InsertBlog(uuid, title, sort, tagList)
	}
	response.JsonResponse(w, framework.ErrorOK)
}

func (f *FileController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		response.JsonResponse(w, framework.ErrorMethodError)
		return
	}

	f.SessionController.HandlerRequest(f, w, r)
	// if status, err := f.WebSession.Get("status"); err != nil || status != "auth" {
	// 	fmt.Println("err: ", err.Error(), "\tstatus: ", status)
	// 	response.JsonResponseWithMsg(w, framework.ErrorAccountAuthError, "not auth")
	// 	return
	// }
	// if r.Header.Get("Content-Type") == "multipart/form-data" {
	// 	fmt.Println("right")
	// }
	switch r.URL.Path {
	case "/personal/file-download":
		f.handlerDownloadRequest(w, r)
	case "/personal/file-upload":
		f.handlerUploadRequest(w, r)
	default:
		response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
	}
}
