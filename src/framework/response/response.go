package response

import (
	"fmt"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, errorCode int) {
	fmt.Fprintf(w, `{"code": %d}`, errorCode)
}

func JsonResponseWithMsg(w http.ResponseWriter, errorCode int, msg string) {
	fmt.Fprintf(w, `{"code": %d, "msg": %s}`, errorCode, msg)
}

func NotFoundResponse() {

}
