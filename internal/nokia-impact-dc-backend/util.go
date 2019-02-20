package nokia_impact_dc_backend

import (
	"cloud.google.com/go/firestore"
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

func UID(r *http.Request) string {
	if tok, ok := r.Context().Value("uid").(string); !ok {
		// BIG screw up if we get here.
		// AuthMiddleware should ensure this always exists.
		panic("no token in context")
		return ""
	} else {
		return tok
	}
}

func GetFirestoreClient(r *http.Request) *firestore.Client {
	if db, ok := r.Context().Value("db").(*firestore.Client); !ok {
		// BIG screw up if we get here.
		panic("no db pool in request context")
		return nil
	} else {
		return db
	}
}

func DB(r *http.Request) *Database {
	return NewDatabase(GetFirestoreClient(r))
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
