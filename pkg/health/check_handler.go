package health

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
)

const talapiHealthPath = "/health"

type healthResponse struct {
	Version string `json:"version"`
	Talapi  string `json:"talapi"`
}

type talapiHealthResponse struct {
	Message string `json:"message"`
}

func talapiStatus(talapiHealthCheckURL url.URL) string {
	log.Println("Calling talapi healthcheck at ", talapiHealthCheckURL.String())

	req, err := http.NewRequest("GET", talapiHealthCheckURL.String(), nil)
	if err != nil {
		return "An error occurred building the Talapi healtcheck request: " + err.Error()
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "An error occurred calling Talapi: " + err.Error()
	}

	if res.StatusCode != http.StatusOK {
		return "Talapi responded with status: " + res.Status
	}

	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var health talapiHealthResponse
	if err := decoder.Decode(&health); err != nil {
		return "Failed to decode talapi healthcheck response: " + err.Error()
	}

	log.Println("Talapi healthcheck: " + health.Message)
	return health.Message
}

// CheckHandler returns the health status of the api
func CheckHandler(talapiBaseURL url.URL, version string) http.Handler {
	talapiHealthCheckURL := talapiBaseURL
	talapiHealthCheckURL.Path = path.Join(talapiBaseURL.Path, talapiHealthPath)

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if err := json.NewEncoder(rw).Encode(healthResponse{
			Talapi:  talapiStatus(talapiHealthCheckURL),
			Version: version,
		}); err == nil {
			rw.Header().Set("Content-Type", "application/json; charset=utf-8")
			rw.WriteHeader(http.StatusOK)
		} else {
			message := "Failed to encode healthcheck result: " + err.Error()
			log.Println(message)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(message))
		}
	})
}
