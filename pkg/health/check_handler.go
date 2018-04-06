package health

import (
	"damo/pkg/config"
	"damo/pkg/util"
	"encoding/json"
	"log"
	"net/http"
)

type healthResponse struct {
	Version config.Version `json:"version"`
}

// CheckHandler returns the health status of the api
func CheckHandler(version config.Version) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		if err := json.NewEncoder(rw).Encode(healthResponse{
			Version: version,
		}); err != nil {
			message := "Failed to encode healthcheck result: " + err.Error()
			log.Println(message)
			util.WriteError(&rw, http.StatusInternalServerError, message)
		}
	})
}
