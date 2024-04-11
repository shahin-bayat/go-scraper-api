package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}, headers map[string]string) error {
	w.Header().Set("Content-Type", "application/json")

	for key, value := range headers {
		_, ok := headers[key]
		if !ok {
			continue
		}
		w.Header().Set(key, value)
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteErrorJSON(w http.ResponseWriter, status int, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorResponse := ErrorResponse{Error: err.Error()}
	return json.NewEncoder(w).Encode(errorResponse)
}

func ReadBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

func DecodeRequestBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}
