package backend

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/base64"
	"encoding/json"
	"firebase.google.com/go"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

var gRegistered = false
var gConfig *Config

func init() {
	NewBackend("ubiik-auth.json")
}

func config() *Config {
	return gConfig
}

func ImpactURL(path string) string {
	return "https://impact.idc.nokia.com" + path
}

type Backend struct {
	DB     *firestore.Client
	Router *mux.Router
}

// like a regular handler but with return codes
type VmeHandlerFunc func(http.ResponseWriter, *http.Request) (int, error)
type VmeHandler struct {
	Handler VmeHandlerFunc
}

func (ah VmeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// make doubly sure we have an auth token
	if _, ok := r.Context().Value("uid").(string); !ok {
		// This would be a big screw up. The other place is db-util.go:Token()... but we can handle it better here.
		log.Println("somehow arrived in handler without uid in the context, auth layer must be broken")
		HttpError(w, http.StatusUnauthorized)
		return
	} else {
		// ValidateToken? If a token exists it should already have been validated in auth.go when it was deserialised
	}

	// Updated to pass ah.appContext as a parameter to our handler type.
	stat, err := ah.Handler(w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", stat, err)
		switch stat {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page, we can
			// now leverage our context instance - e.g.
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(stat), stat)
		default:
			http.Error(w, http.StatusText(stat), stat)
		}
	}
}

func NewBackend(authFile string) *Backend {
	//fmt.Println("Loading gConfig...")
	gConfig = ReadConfig()

	//fmt.Println("Initialising db...")
	ctx := context.Background()

	// Use a service account
	sa := option.WithCredentialsFile(authFile)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println("Initialising router...")
	router := mux.NewRouter()

	router.HandleFunc("/test", TestEndpoint).Methods("GET")
	router.HandleFunc("/_ah/warmup", WarmupEndpoint).Methods("GET")

	// inject various bits and pieces (that would otherwise be global) into context
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", firestoreClient)
			ctx = context.WithValue(ctx, "authclient", authClient)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	v1Router := router.PathPrefix("/api/impact/v1").Subrouter()
	v1Router.Use(AuthMiddleware)
	v1Router.Handle("/callback", VmeHandler{CallbackEndpoint}).Methods("POST")



	//v1Router.Handle("/group/{groupId}", VmeHandler{GroupUpdateEndpoint}).Methods("PUT")

	http.Handle("/", router)

	return &Backend{DB: firestoreClient, Router: router}
}


//func (nokia-impact-dc *Backend) Server(address string) {
//
//	fmt.Printf("Starting server at http://%s...\n", address)
//
//	log.Fatal(
//		http.ListenAndServe(address,
//			handlers.CORS(
//				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
//				// TODO: tighten up CORS
//				handlers.AllowedOrigins([]string{"*"}))(nokia-impact-dc.Router)))
//}

func transformError(err error) (int, error) {
	code := status.Code(err)
	switch code {
	case codes.OK:
		return OK()
	case codes.AlreadyExists:
		return Conflict(err)
	case codes.NotFound:
		return NotFound(err)
	case codes.PermissionDenied, codes.Unauthenticated:
		return Unauthorized(err)
	default:
		return InternalError(err)
	}
}
func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Your request: %+v", r)
}

// register with the Nokia IMPACT platform


func WarmupEndpoint(w http.ResponseWriter, r *http.Request) {
	type RegReqHdr struct {
		Authorization string `json:"authorization"`
	}
	type RegReq struct {
		Headers RegReqHdr `json:"headers"`
		Url string `json:"url"`
	}
	authToken := gConfig.CallbackUsername + ":" + gConfig.CallbackPassword
	authToken = "Basic " + base64.StdEncoding.EncodeToString([]byte(authToken))
	data := &RegReq{
		Headers: RegReqHdr{ authToken},
		Url:     gConfig.CallbackHost + "/api/impact/v1/callback",
	};
	buf, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(string(buf))
	rdr := bytes.NewReader(buf)

	url := ImpactURL("/m2m/applications/registration")
	rv := HttpPut(r, url, rdr)
	
	log.Printf("rv %+v\n", rv)
	gRegistered = true
}

func CallbackEndpoint(w http.ResponseWriter, r *http.Request) (int, error) {
	//db := DB(r)
	//rv, err := db.GetAllMeta(UID(r))
	//if err != nil {
	//	return InternalError(err)
	//}
	//json.NewEncoder(w).Encode(rv)
	return OK()
}

//func GroupUpdateEndpoint(w http.ResponseWriter, r *http.Request) (int, error) {
//	vars := mux.Vars(r)
//	groupId, ok := vars["groupId"]
//	if !ok {
//		return BadRequest(errors.New("no group id on path"))
//	}
//	group := structs.Group{}
//	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
//		return BadRequest(err)
//	}
//	if groupId != group.Id {
//		return BadRequest(errors.New("mismatched group id"))
//	}
//	db := DB(r)
//	rv, err := db.UpdateGroup(UID(r), &group)
//	if err != nil {
//		return transformError(err)
//	}
//	json.NewEncoder(w).Encode(rv)
//	return OK()
//}
