package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	StatusCode int `json:"status_code"`
	Message    any `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{StatusCode: statusCode, Message: err.Error()}
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{StatusCode: http.StatusBadRequest, Message: errors}
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, errors.New("invalid JSON request data"))
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
		return err
	}
	return nil
}

func Make(h APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			var apiErr APIError
			if errors.As(err, &apiErr) {
				WriteJSON(w, apiErr.StatusCode, apiErr.Message, nil)
			} else {
				errResp := NewAPIError(http.StatusInternalServerError, errors.New("internal server error"))
				WriteJSON(w, errResp.StatusCode, errResp.Message, nil)
			}
			slog.Error(
				"HTTP API error", "error", err.Error(), "path", r.URL.Path,
			)

		}
	}
}
