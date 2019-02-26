package nokia_impact_dc_backend

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ConfigType struct {
	CallbackUsername string
	CallbackPassword string
	ImpactUsername   string
	ImpactPassword   string
	ImpactBaseURL    string
	ImpactGroup      string
	GoogleAuthFile   string
	ListenPort       string
}

func LoadConfig() *ConfigType {
	fp, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Reading config file", err)
	}
	cfg := ConfigType{}
	err = json.Unmarshal(fp, &cfg)
	if err != nil {
		log.Fatalln("Parsing config file", err)
	}
	return &cfg
}
