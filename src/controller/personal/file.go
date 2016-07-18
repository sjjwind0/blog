package personal

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type FileController struct {
}

func NewFileController() *FileController {
	return &FileController{}
}

func (f *FileController) Path() interface{} {
	return "/file"
}

func (s *FileController) HandlerRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("/tmp/test.upload", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}
