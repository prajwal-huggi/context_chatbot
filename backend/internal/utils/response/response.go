package response

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status_code int, data interface{}) error{
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)

	return json.NewEncoder(w).Encode(data)
}