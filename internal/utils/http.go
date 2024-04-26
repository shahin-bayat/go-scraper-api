package utils

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}, headers map[string]string) {
	w.Header().Set("Content-Type", "application/json")

	for key, value := range headers {
		_, ok := headers[key]
		if !ok {
			continue
		}
		w.Header().Set(key, value)
	}
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Fatal(err)
	}
}

func WriteErrorJSON(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	errorResponse := ErrorResponse{Error: err.Error()}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		log.Fatal(err)
	}
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
	err := decoder.Decode(v)
	if err != nil {
		return errors.New("failed to decode request body")
	}
	return nil
}

func DecodeResponseBody(body io.ReadCloser, v interface{}) error {
	err := json.NewDecoder(body).Decode(v)
	if err != nil {
		return errors.New("failed to decode response body")
	}
	return nil
}
