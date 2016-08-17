package response

import (
	"fmt"
	"framework/base/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, errorCode int) {
	fmt.Fprintf(w, `{"code": %d}`, errorCode)
}

func JsonResponseWithMsg(w http.ResponseWriter, errorCode int, msg string) {
	fmt.Fprintf(w, `{"code": %d, "msg": "%s"}`, errorCode, msg)
}

func JsonResponseWithData(w http.ResponseWriter, errorCode int, msg string, data interface{}) {
	fmt.Fprintf(w, `{"code": %d, "msg": "%s", "data": %s}`, errorCode, msg, json.ToJsonString(data))
}

func NotFoundResponse() {

}
