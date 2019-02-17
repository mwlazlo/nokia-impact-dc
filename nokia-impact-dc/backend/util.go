package backend

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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


func HttpPut(r *http.Request, url string, rdr io.Reader) string {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	request, err := http.NewRequest("PUT", url, rdr)
	request.SetBasicAuth(gConfig.ImpactUsername, gConfig.ImpactPassword)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("The calculated length is:", len(string(contents)), "for the url:", url)
		log.Println("   ", response.StatusCode)
		hdr := response.Header
		for key, value := range hdr {
			log.Println("   ", key, ":", value)
		}
		rv := string(contents)
		log.Println(rv)
		return rv
	}
	return ""
}