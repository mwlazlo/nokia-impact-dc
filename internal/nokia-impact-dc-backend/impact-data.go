package nokia_impact_dc_backend

import (
	"encoding/json"
	"fmt"
	"log"
)

type impactData struct {
	Reports         interface{} `json:"reports"`
	Registrations   interface{} `json:"registrations"`
	Deregistrations interface{} `json:"deregistrations"`
	Updates         interface{} `json:"updated"`
	Expirations     interface{} `json:"expirations"`
	Responses       interface{} `json:"responses"`
}

type AbstractDataRecord struct {
	UpdateType     string // registation / deregistration, etc
	SubscriptionID string
	SerialNumber   string
	IMSI           string
	Timestamp      int64

	NumberValues  map[string]float64
	StringValues  map[string]string
	BooleanValues map[string]bool
	ArrayValues   map[string][]interface{}
}

func newAbstractDataRecord(updateType string) *AbstractDataRecord {
	abs := AbstractDataRecord{}
	abs.UpdateType = updateType
	abs.NumberValues = make(map[string]float64)
	abs.StringValues = make(map[string]string)
	abs.BooleanValues = make(map[string]bool)
	abs.ArrayValues = make(map[string][]interface{})
	return &abs
}

func ParseImpactJSON(data []byte) []*AbstractDataRecord {
	rv := make([]*AbstractDataRecord, 0)
	var f impactData
	err := json.Unmarshal(data, &f)
	if err != nil {
		log.Fatal("Failed to parse JSON")
	}

	process := func(i interface{}, updateType string) {
		if i != nil {
			for _, v := range (i.([]interface{})) {
				x := v.(map[string]interface{})
				abs := newAbstractDataRecord(updateType)
				parseStruct(x, abs)
				rv = append(rv, abs)
			}
		}
	}

	process(f.Registrations, "registrations")
	process(f.Deregistrations, "deregistrations")
	process(f.Expirations, "expirations")
	process(f.Reports, "reports")
	process(f.Responses, "responses")
	process(f.Updates, "updates")

	return rv
}

func parseStruct(m map[string]interface{}, abs *AbstractDataRecord) {
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			switch k {
			case "imsi":
				abs.IMSI = vv
			case "serialNumber":
				abs.SerialNumber = vv
			case "subscriptionId":
				abs.SubscriptionID = vv
			}
			abs.StringValues[k] = vv
		case float64:
			if k == "timestamp" {
				abs.Timestamp = int64(vv)
			}
			abs.NumberValues[k] = vv
		case bool:
			abs.BooleanValues[k] = vv
		case []interface{}:
			abs.ArrayValues[k] = vv
			//fmt.Println(prefix, k, "is an array:")
			//for i, u := range vv {
			//	fmt.Println(prefix, i, u)
			//}
		case map[string]interface{}:
			parseStruct(vv, abs)
		default:
			fmt.Println(k, "is a type we can't handle", vv)
		}
	}
}
