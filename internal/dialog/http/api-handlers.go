package internalhttp

import (
	"encoding/json"
	"errors"
	dialogapp "github.com/filatkinen/socialnet/internal/dialog/app"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *DialogService) LoginPost(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	var login LoginBody
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		s.ClientError(w, http.StatusBadRequest, "Wrong JSON request")
		return
	}
	// TODO s.app.UserLogin - change to GRPC
	token, err := s.UserLogin(r.Context(), login.Id, login.Password)
	if err != nil {
		if errors.Is(err, dialogapp.ErrorUserNotFound) {
			s.ClientError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, dialogapp.ErrorUserPassInvalid) {
			s.ClientError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: rID,
			Code:      500,
		})
		return
	}

	var ulogin UserLoginResponse
	ulogin.Token = token
	s.writeHTTPJsonOK(w, ulogin)
}

func (s *DialogService) DialogUserIdListGet(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	this_userID := r.Context().Value(ContextUserKey).(string)

	vars := mux.Vars(r)
	userid, ok := vars["user_id"]
	if !ok {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   "user id was not set in URL",
			RequestId: rID,
			Code:      500,
		})
		return
	}
	messages, err := s.app.UserDialogListdMessages(r.Context(), this_userID, userid)
	if err != nil {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: rID,
			Code:      500,
		})
		return
	}
	s.writeHTTPJsonOKIdent(w, messages)
}

func (s *DialogService) DialogUserIdSendPost(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	this_userID := r.Context().Value(ContextUserKey).(string)

	vars := mux.Vars(r)
	userid, ok := vars["user_id"]
	if !ok {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   "user id was not set in URL",
			RequestId: rID,
			Code:      500,
		})
		return
	}

	data := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: rID,
			Code:      500,
		})
		return
	}
	text, ok := data["text"].(string)
	if !ok {
		s.ClientError(w, http.StatusBadRequest, "")
		return
	}
	err = s.app.UserDialogSendMessage(r.Context(), this_userID, userid, text)
	if err != nil {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: rID,
			Code:      500,
		})
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
