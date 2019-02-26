package nokia_impact_dc_backend

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
)

func InternalError(err error) (int, error) {
	return http.StatusInternalServerError, err
}

func Conflict(err error) (int, error) {
	return http.StatusConflict, err
}

func NotFound(err error) (int, error) {
	return http.StatusNotFound, err
}

func BadRequest(err error) (int, error) {
	return http.StatusBadRequest, err
}

func OK() (int, error) {
	return http.StatusOK, nil
}

func Unauthorized(err error) (int, error) {
	return http.StatusUnauthorized, err
}

func HttpError(w http.ResponseWriter, code int) {
	log.Printf("HttpError(%d)\n", code)
	http.Error(w, http.StatusText(code), code)
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
