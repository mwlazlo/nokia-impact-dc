package nokia_impact_dc_backend

import (
	"context"
	"firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

// For now, just accept all data. Later on, filter based on what subscriptions we're interested in.
// Search for this tag in related comments: SUBSC
type AppContext struct {
	Db     *Database
	Config *ConfigType
	//ResourceSubscriptions map[string]bool
	//LifecycleSubscriptions map[string]bool
}

// SUBSC
//func (ctx *AppContext) addResourceSubscriptionId(id string) {
//	mutex := sync.Mutex{}
//	defer mutex.Unlock()
//	mutex.Lock()
//	ctx.ResourceSubscriptions[id] = true
//}
//
//func (ctx *AppContext) addLifecycleSubscriptionId(id string) {
//	mutex := sync.Mutex{}
//	defer mutex.Unlock()
//	mutex.Lock()
//	ctx.LifecycleSubscriptions[id] = true
//}
//
func (ctx *AppContext) hasResourceSubscriptionId(id string) bool {
	//	_, ok := ctx.ResourceSubscriptions[id]
	//	return ok
	return true // SUBSC: accept everything for now
}

//
func (ctx *AppContext) hasLifecycleSubscriptionId(id string) bool {
	//	_, ok := ctx.LifecycleSubscriptions[id]
	//	return ok
	return true // SUBSC: accept everything for now
}

// http.Handler wrapper to provide return codes and context
type AppHandlerFunc func(*AppContext, http.ResponseWriter, *http.Request) (int, error)

type AppHandler struct {
	*AppContext
	Handler AppHandlerFunc
}

func (ah AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// call the handler function and do any special handling of return codes
	log.Println("Handler:", GetFunctionName(ah.Handler))
	status, err := ah.Handler(ah.AppContext, w, r)
	if err != nil {
		log.Printf("HTTP %d: %q", status, err)
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
			// And if we wanted a friendlier error page, we can
			// now leverage our context instance - e.g.
			// err := ah.renderTemplate(w, "http_404.tmpl", nil)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			http.Error(w, http.StatusText(status), status)
		}
	}
}

// main entry point
func Backend() {
	log.Println("Initialising db...")

	ctx := context.Background()
	cfg := LoadConfig()

	// Use a service account
	sa := option.WithCredentialsFile(cfg.GoogleAuthFile)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	appCtx := &AppContext{
		NewDatabase(firestoreClient),
		cfg,
		// SUBSC
		//make(map[string]bool),
		//make(map[string]bool),
	}

	log.Println("Initialising router...")
	router := mux.NewRouter()

	router.Handle("/test", AppHandler{appCtx, TestEndpoint})

	v1Router := router.PathPrefix("/v1").Subrouter()
	v1Router.Use(AuthMiddleware(appCtx))
	v1Router.Handle("/callback", AppHandler{appCtx, CallbackEndpoint}).Methods("POST")

	// SUBSC
	//v1Router.Handle("/addLifecycleSubscription", AppHandler{appCtx,SetLifecycleSubEndpoint}).Methods("POST")
	//v1Router.Handle("/addResourceSubscription", AppHandler{appCtx,SetResourceSubEndpoint}).Methods("POST")

	log.Println("Finished initialisation")
	log.Println("Listening on port", cfg.ListenPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ListenPort, router))
}
