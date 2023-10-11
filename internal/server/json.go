package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.MarshalIndent(payload, "", "	")
	if err != nil {
		log.Printf("Failed to marshal payload: %v", payload)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
	log.Printf("\nResponse Code: %v \nResponse Data:\n%v", code, string(data))
}
