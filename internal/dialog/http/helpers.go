package internalhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *DialogService) ClientError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if message == "" {
		fmt.Fprintln(w, http.StatusText(status))
		return
	}
	fmt.Fprintf(w, "%s: %s\n", http.StatusText(status), message)
}

func (s *DialogService) ServerError(w http.ResponseWriter, status int, resp *InlineResponse500) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Retry-After", "10")
	w.WriteHeader(status)
	value, err := s.toJSON(resp)
	if err != nil {
		fmt.Fprintln(w, "error marshaling response with mistake details")
	}
	fmt.Fprintln(w, value)
}

func (s *DialogService) toJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *DialogService) toJSONIdent(value any) (string, error) {
	data, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *DialogService) writeHTTPJsonOK(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data, err := s.toJSON(value)
	if err != nil {
		fmt.Fprintln(w, "request OK")
		fmt.Fprintln(w, "error marshaling response with details of request")
		return
	}
	fmt.Fprintln(w, data)
}

func (s *DialogService) writeHTTPJsonOKIdent(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	data, err := s.toJSONIdent(value)
	if err != nil {
		fmt.Fprintln(w, "request OK")
		fmt.Fprintln(w, "error marshaling response with details of request")
		return
	}
	fmt.Fprintln(w, data)
}
