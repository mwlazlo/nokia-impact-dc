package main

import (
	"log"
	"net/http"
	"nokia-impact-dc/internal/nokia-impact-dc-backend"
)

func main() {
	nokia_impact_dc_backend.InitConfig()
	nokia_impact_dc_backend.InitialiseBackend()
	log.Println("Listening on port", nokia_impact_dc_backend.Config().ListenPort)
	log.Fatal(http.ListenAndServe(":"+nokia_impact_dc_backend.Config().ListenPort, nil))
}
