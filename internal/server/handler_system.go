package server

import (
	"net/http"
)

func CheckReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, map[string]string{"ready": "OK"})
}
