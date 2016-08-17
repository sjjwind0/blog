package personal

import (
	"fmt"
	"framework"
	"framework/base/archive"
	"framework/base/config"
	"framework/base/json"
	"framework/response"
	"framework/server"
	"io"
	"model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const k24K = (1 << 20) * 24

type FileController struct {
	server.SessionController
}

func NewPersonalFileController() *FileController {
	return &FileController{}
}

func (f *FileController) Path() interface{} {
	return "/personal/file"
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
	rawPath := config.GetDefaultConfigJsonReader().Get("blog.storage.file.raw").(string)
	blogInfo, err := model.ShareBlogModel().FetchBlogByBlogID(blogId)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
		return
	}
	blogPath := filepath.Join(rawPath, blogInfo.BlogUUID+".zip")
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
	w.Header().Set("Content-Disposition", "attachment; filename="+blogInfo.BlogUUID+".zip")
	w.Write(content)
}

func (f *FileController) savePostFile(r *http.Request, name string, path string) string {
	file, handler, err := r.FormFile(name)
	if err != nil {
		fmt.Println("r.FormFile: ", err)
		return ""
	}
	defer file.Close()
	saveFile, err := os.OpenFile(filepath.Join(path, handler.Filename), os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("os.OpenFile: ", err)
		return ""
	}
	defer saveFile.Close()
	io.Copy(saveFile, file)
	return handler.Filename
}

func (f *FileController) checkFolder(path string) {
	_, err := os.Stat(path)
	if !os.IsExist(err) {
		fmt.Println("create folder: ", path)
		err := os.MkdirAll(path, 0775)
		if err != nil {
			fmt.Println("create folder error: ", err.Error())
		}
	}
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
**	 			- other
** raw文件放在raw目录不对外开放，html文件以及res文件放在blog目录，meta信息放数据库。
 */
func (f *FileController) handlerUploadRequest(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(k24K); nil != err {
		fmt.Println("r.ParseMultipartForm: ", err)
		return
	}

	saveTmpPath := "/tmp"
	rawZipName := f.savePostFile(r, "raw", saveTmpPath)
	webHtmlName := f.savePostFile(r, "web", saveTmpPath)
	blogInfoName := f.savePostFile(r, "info", saveTmpPath)
	resZipName := f.savePostFile(r, "res", saveTmpPath)
	coverImgName := f.savePostFile(r, "img", saveTmpPath)

	blogMetaInfoReader := json.NewJsonReader(filepath.Join(saveTmpPath, blogInfoName))
	uuid := blogMetaInfoReader.Get("uuid").(string)
	title := blogMetaInfoReader.Get("title").(string)
	tag := blogMetaInfoReader.Get("tag").(string)
	tagList := strings.Split(tag, "||")
	sort := blogMetaInfoReader.Get("sort").(string)
	isExist, err := model.ShareBlogModel().BlogIsExistByUUID(uuid)
	if err != nil {
		response.JsonResponseWithMsg(w, framework.ErrorSQLError, err.Error())
		fmt.Println("insert uuid error: ", err.Error())
		return
	}
	// 7. archive to path
	rawRootPath := config.GetDefaultConfigJsonReader().Get("blog.storage.file.raw").(string)
	f.checkFolder(rawRootPath)
	blogRootPath := config.GetDefaultConfigJsonReader().Get("blog.storage.file.blog").(string)
	blogRootPath = filepath.Join(blogRootPath, uuid)
	f.checkFolder(blogRootPath)

	rawZipPath := filepath.Join(rawRootPath, uuid+".zip")
	fmt.Println("1 = ", filepath.Join(saveTmpPath, rawZipName))
	fmt.Println("2 = ", rawZipPath)
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
		fmt.Println("unzip error: ", err.Error())
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
	f.SessionController.HandlerRequest(f, w, r)
	// if status, err := f.WebSession.Get("status"); err != nil || status != "auth" {
	// 	fmt.Println("err: ", err.Error(), "\tstatus: ", status)
	// 	response.JsonResponseWithMsg(w, framework.ErrorAccountAuthError, "not auth")
	// 	return
	// }
	// if r.Header.Get("Content-Type") == "multipart/form-data" {
	// 	fmt.Println("right")
	// }
	fmt.Println("FileController.HandlerRequest")
	switch r.Method {
	case "POST":
		f.handlerUploadRequest(w, r)
	case "GET":
		f.handlerDownloadRequest(w, r)
	default:
		response.JsonResponseWithMsg(w, framework.ErrorParamError, "param error")
	}
}
