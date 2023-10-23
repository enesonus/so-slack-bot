package server

import (
	"net/http"
	"os"
	"strings"
)

func CheckReadiness(w http.ResponseWriter, r *http.Request) {
	message := "Server"
	hostname, _ := os.Hostname()
	if strings.HasSuffix(hostname, "local") {
		message = "Localhost"
	}
	respondWithJSON(w, 200, map[string]string{"ready": message})
}
