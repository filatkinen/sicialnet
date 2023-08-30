package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) ClientError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if message == "" {
		fmt.Fprintln(w, http.StatusText(status))
		return
	}
	fmt.Fprintf(w, "%s: %s\n", http.StatusText(status), message)
}

func (s *Server) ServerError(w http.ResponseWriter, status int, resp *InlineResponse500) {
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

func (s *Server) toJSON(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) toJSONIdent(value any) (string, error) {
	data, err := json.MarshalIndent(value, "", " ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Server) writeHTTPJsonOK(w http.ResponseWriter, value any) {
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

func (s *Server) writeHTTPJsonOKIdent(w http.ResponseWriter, value any) {
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

func (s *Server) writeHTTPTextOK(w http.ResponseWriter, value string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, value)
}

func (s *Server) PutPostMessageToRabbit(postID string, postText string, userID string) error {
	friends, err := s.app.UserGetFriends(context.Background(), userID)
	if err != nil {
		return err
	}
	post := Post{
		Id:           postID,
		Text:         postText,
		AuthorUserId: userID,
	}
	b, err := json.Marshal(post)
	if err != nil {
		return err
	}
	go func() {
		for _, val := range friends {
			s.rabbit.SendMessages(b, val)
		}
	}()
	return nil
}
