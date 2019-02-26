package nokia_impact_dc_backend

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type SubUpdate struct {
	Id string `json:"id"`
}

func TestEndpoint(ctx *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Fprintf(w, "ConfigType = %+v\n", ctx.Config)
	fmt.Fprintf(w, "Your request: %+v", r)
	return OK()
}

func CallbackEndpoint(ctx *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
	log.Printf("Callback request: %+v\n", r)

	jsonObject, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("Could not read body", err)
	}

	if len(jsonObject) == 0 {
		return OK()
	}

	log.Println("Body", string(jsonObject))
	objects := ParseImpactJSON(jsonObject)

	for _, v := range objects {
		if ctx.hasLifecycleSubscriptionId(v.SubscriptionID) || ctx.hasResourceSubscriptionId(v.SubscriptionID) {
			if err := ctx.Db.saveData(r, v); err != nil {
				return InternalError(err)
			}
		} else {
			log.Println("Ignoring subscription", v.SubscriptionID)
		}
	}

	return OK()
}

// SUBSC
//func SetLifecycleSubEndpoint(ctx *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
//	upd := SubUpdate{}
//	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
//		return BadRequest(err)
//	}
//	log.Println("Updated lifecycle event subscription", upd.Id)
//	ctx.addLifecycleSubscriptionId(upd.Id)
//	return OK()
//}
//
//func SetResourceSubEndpoint(ctx *AppContext, w http.ResponseWriter, r *http.Request) (int, error) {
//	upd := SubUpdate{}
//	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
//		return BadRequest(err)
//	}
//	log.Println("Updated resource event subscription", upd.Id)
//	ctx.addResourceSubscriptionId(upd.Id)
//	return OK()
//}
