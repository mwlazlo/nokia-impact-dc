package nokia_impact_dc_backend

import (
	"context"
	"firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

// like a regular handler but with return codes
type VmeHandlerFunc func(http.ResponseWriter, *http.Request) (int, error)
type VmeHandler struct {
	Handler VmeHandlerFunc
}

// wrap an http handler. The main reason is to ensure the request is authenticated.
func (ah VmeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// make doubly sure we have an auth token
	if _, ok := r.Context().Value("uid").(string); !ok {
		log.Println("somehow arrived in handler without uid in the context, auth layer must be broken?")
		HttpError(w, http.StatusUnauthorized)
		return
	} else {
		// ValidateToken? If a token exists it should already have been validated in auth.go when it was deserialised
	}

	// call the handler function and do any special handling of return codes
	log.Println("Handler:", GetFunctionName(ah.Handler))
	stat, err := ah.Handler(w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", stat, err)
		switch stat {
		case http.StatusOK:
			// noop
		case http.StatusNotFound:
			http.NotFound(w, r)
		default:
			http.Error(w, http.StatusText(stat), stat)
		}
	}
}

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

func InitialiseBackend() {
	log.Println("Initialising db...")

	ctx := context.Background()

	// Use a service account
	sa := option.WithCredentialsFile(Config().GoogleAuthFile)
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

	log.Println("Initialising router...")
	router := mux.NewRouter()

	router.HandleFunc("/test", TestEndpoint)

	// inject various bits and pieces (that would otherwise be global) into context
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", firestoreClient)
			ctx = context.WithValue(ctx, "authclient", authClient)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	v1Router := router.PathPrefix("/v1").Subrouter()
	v1Router.Use(AuthMiddleware)
	v1Router.Handle("/callback", VmeHandler{CallbackEndpoint}).Methods("POST")
	v1Router.Handle("/setLifecycleSubscription", VmeHandler{SetLifecycleSubEndpoint}).Methods("POST")
	v1Router.Handle("/setResourceSubscription", VmeHandler{SetResourceSubEndpoint}).Methods("POST")

	http.Handle("/", router)

	log.Println("Finished initialisation", Config())
}
