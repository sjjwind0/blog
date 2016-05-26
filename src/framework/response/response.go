package response

import (
	"fmt"
	"net/http"
)

const (
	ErrorOK         = iota
	ErrorParamError = iota
)

func JsonResponse(w http.ResponseWriter, errorCode int) {
	fmt.Fprintf(w, "{code: %d}", errorCode)
}

func JsonResponseWithMsg(w http.ResponseWriter, errorCode int, msg string) {
	fmt.Fprintf(w, "{code: %d, msg: %s}", errorCode, msg)
}

func NotFoundResponse() {

}
