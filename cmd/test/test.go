package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type ImpactData struct {
	Reports interface{}			`json:"reports"`
	Registrations interface{} `json:"registrations"`
	Deregistrations interface{} `json:"deregistrations"`
	Updates interface{} `json:"updated"`
	Expirations interface{} `json:"expirations"`
	Responses interface{} `json:"responses"`
}

type AbstractDataRecord struct {
	UpdateType string // registation / deregistration
	SubscriptionID string
	SerialNumber string
	IMSI string

	NumberValues map[string]float64
	StringValues map[string]string
	BooleanValues map[string]bool
	ArrayValues map[string][]interface{}
}

func NewAbstractDataRecord() *AbstractDataRecord {
	abs := AbstractDataRecord{}
	abs.NumberValues = make(map[string]float64)
	abs.StringValues = make(map[string]string)
	abs.BooleanValues = make(map[string]bool)
	abs.ArrayValues = make(map[string][]interface{})
	return &abs
}

var DATA = "{\"reports\":[],\"registrations\":[{\"deviceType\":\"device\",\"serialNumber\":\"c23acba8-34b5-11e9-ad21-aba0f2afd3df\",\"timestamp\":1550629113761,\"make\":\"Open Mobile Alliance\",\"model\":\"Lightweight M2M Client\",\"firmwareVersion\":\"1.0\",\"groupName\":\"APJ.JAPAN.Rakuten\",\"imsi\":\"N/A\",\"address\":\"345000123\",\"protocol\":\"LWM2M\",\"tags\":\"N/A\",\"subscriptionId\":\"32394e5a-ed85-4bd0-bb81-8151a2bb7619\",\"deviceProps\":{\"device/0/manufacturer\":\"Open Mobile Alliance\",\"device/0/model\":\"Lightweight M2M Client\",\"device/0/serialNumber\":\"345000123\",\"device/0/firmwareVersion\":\"1.0\",\"device/0/availablePowerSources/0\":\"1\",\"device/0/availablePowerSources/1\":\"5\",\"device/0/powerSourceVoltage/0\":\"3800\",\"device/0/powerSourceVoltage/1\":\"5000\",\"device/0/powerSourceCurrent/0\":\"125\",\"device/0/powerSourceCurrent/1\":\"900\",\"device/0/batteryLevel\":\"100\",\"device/0/freeMemory\":\"15\",\"device/0/errorCode/0\":\"0\",\"device/0/currentTime\":3101258174000,\"device/0/utcOffset\":\"+01:00\",\"device/0/timezone\":\"Europe/Berlin\",\"device/0/currentBindingMode\":\"U\",\"location/0/latitude\":27.986064910888672,\"location/0/longitude\":86.92262268066406,\"location/0/altitude\":8495.0,\"location/0/radius\":0.0,\"location/0/velocity\":\"\",\"location/0/timestamp\":1550629086000,\"location/0/speed\":0.0,\"device/0/endPointClientName\":\"c23acba8-34b5-11e9-ad21-aba0f2afd3df\"}},{\"deviceType\":\"device\",\"serialNumber\":\"c23acba8-34b5-11e9-ad21-aba0f2afd3df\",\"timestamp\":1550629113761,\"make\":\"Open Mobile Alliance\",\"model\":\"Lightweight M2M Client\",\"firmwareVersion\":\"1.0\",\"groupName\":\"APJ.JAPAN.Rakuten\",\"imsi\":\"N/A\",\"address\":\"345000123\",\"protocol\":\"LWM2M\",\"tags\":\"N/A\",\"subscriptionId\":\"b33d4565-2f62-4b15-88a5-ad41ffaff284\",\"deviceProps\":{\"device/0/manufacturer\":\"Open Mobile Alliance\",\"device/0/model\":\"Lightweight M2M Client\",\"device/0/serialNumber\":\"345000123\",\"device/0/firmwareVersion\":\"1.0\",\"device/0/availablePowerSources/0\":\"1\",\"device/0/availablePowerSources/1\":\"5\",\"device/0/powerSourceVoltage/0\":\"3800\",\"device/0/powerSourceVoltage/1\":\"5000\",\"device/0/powerSourceCurrent/0\":\"125\",\"device/0/powerSourceCurrent/1\":\"900\",\"device/0/batteryLevel\":\"100\",\"device/0/freeMemory\":\"15\",\"device/0/errorCode/0\":\"0\",\"device/0/currentTime\":3101258174000,\"device/0/utcOffset\":\"+01:00\",\"device/0/timezone\":\"Europe/Berlin\",\"device/0/currentBindingMode\":\"U\",\"location/0/latitude\":27.986064910888672,\"location/0/longitude\":86.92262268066406,\"location/0/altitude\":8495.0,\"location/0/radius\":0.0,\"location/0/velocity\":\"\",\"location/0/timestamp\":1550629086000,\"location/0/speed\":0.0,\"device/0/endPointClientName\":\"c23acba8-34b5-11e9-ad21-aba0f2afd3df\"}}],\"deregistrations\":[],\"updates\":[],\"expirations\":[],\"responses\":[]}"

func main() {
	var f ImpactData
	b := []byte(DATA)
	err := json.Unmarshal(b, &f)
	if err != nil {
		log.Fatal("Failed to parse JSON")
	}

	m := f.Registrations.([]interface{})

	for _, v := range(m) {
		x := v.(map[string]interface{})
		abs := NewAbstractDataRecord()
		parseStruct(x, abs)
		fmt.Println(abs)
		fmt.Println("-------\n")
	}
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
			default:
				abs.StringValues[k] = vv
			}
		case float64:
			abs.NumberValues[k] = vv
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
