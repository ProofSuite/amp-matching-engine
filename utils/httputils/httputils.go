package httputils

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, code int, message string) {
	WriteJSON(w, code, map[string]string{"error": message})
}

func WriteJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.WriteHeader(code)
	w.Write(response)
}
