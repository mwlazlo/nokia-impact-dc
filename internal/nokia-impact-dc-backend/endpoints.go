package nokia_impact_dc_backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// TODO: find a way that doesn't involve globals. Maybe inject something using middleware
var gLifecycleSubscription = ""
var gResourceSubscription = ""

type SubUpdate struct {
	Id string `json:"id"`
}

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ConfigType = %+v\n", Config())
	fmt.Fprintf(w, "Your request: %+v", r)
}

func CallbackEndpoint(w http.ResponseWriter, r *http.Request) (int, error) {
	log.Printf("Callback request: %+v\n", r)

	jsonObject, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln("Could not read body", err)
	}

	if len(jsonObject) == 0 {
		return OK()
	}

	log.Println("Body", string(jsonObject))

	// TODO pass request and stream body
	objects := ParseImpactJSON(jsonObject)

	for _, v := range objects {
		if v.SubscriptionID == gLifecycleSubscription || v.SubscriptionID == gResourceSubscription {
			if err := DB(r).saveData(r, v); err != nil {
				return InternalError(err)
			}
		} else {
			log.Println("Ignoring subscription", v.SubscriptionID)
		}
	}

	return OK()
}

func SetLifecycleSubEndpoint(w http.ResponseWriter, r *http.Request) (int, error) {
	upd := SubUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		return BadRequest(err)
	}
	log.Println("Updated lifecycle event subscription", upd.Id)
	gLifecycleSubscription = upd.Id
	return OK()
}

func SetResourceSubEndpoint(w http.ResponseWriter, r *http.Request) (int, error) {
	upd := SubUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		return BadRequest(err)
	}
	log.Println("Updated resource event subscription", upd.Id)
	gResourceSubscription = upd.Id
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
