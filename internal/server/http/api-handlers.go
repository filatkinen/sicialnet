package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	socialapp "github.com/filatkinen/socialnet/internal/server/app"
	"github.com/filatkinen/socialnet/internal/storage"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const DateLayout = "2006-01-02"

func (s *Server) LoginPost(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	var login LoginBody
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		s.ClientError(w, http.StatusBadRequest, "Wrong JSON request")
		return
	}
	token, err := s.app.UserLogin(r.Context(), login.Id, login.Password)
	if err != nil {
		if errors.Is(err, socialapp.ErrorUserNotFound) {
			s.ClientError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, socialapp.ErrorUserPassInvalid) {
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

func (s *Server) UserRegisterPost(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	var userReg UserRegisterBody
	err := json.NewDecoder(r.Body).Decode(&userReg)
	if err != nil {
		s.ClientError(w, http.StatusBadRequest, "Wrong JSON request")
		return
	}

	var user storage.User
	if userReg.Birthdate != nil {
		t, err := time.Parse(DateLayout, *userReg.Birthdate)
		if err != nil {
			s.ClientError(w, http.StatusBadRequest, "Wrong date format. Use:YYYY-MM-DD")
			return
		}
		user.BirthDate = &t
	}
	user.City = userReg.City
	user.Biography = userReg.Biography
	user.FirstName = userReg.FirstName
	user.SecondName = userReg.SecondName

	if len(userReg.Password) == 0 {
		s.ClientError(w, http.StatusBadRequest, "Empty password")
		return
	}
	if len(userReg.FirstName) == 0 {
		s.ClientError(w, http.StatusBadRequest, "Empty FirstName")
		return
	}

	userID, err := s.app.UserAdd(r.Context(), &user, userReg.Password)
	if err != nil {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: rID,
			Code:      500,
		})
		return
	}

	var userCreateResponse UserCreateResponse
	userCreateResponse.UserId = userID
	s.writeHTTPJsonOK(w, userCreateResponse)
}

func (s *Server) UserGetIdGet(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	vars := mux.Vars(r)
	userid, ok := vars["id"]
	if !ok {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   "user id was not set in URL",
			RequestId: rID,
			Code:      500,
		})
		return
	}
	u, err := s.app.UserGet(r.Context(), userid)
	if err != nil {
		if errors.Is(err, socialapp.ErrorUserNotFound) {
			s.ClientError(w, http.StatusNotFound, err.Error())
			return
		}
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: rID,
			Code:      500,
		})
		return
	}
	var userGet User
	userGet.FirstName = u.FirstName
	userGet.City = u.City
	userGet.Id = u.Id
	userGet.SecondName = u.SecondName
	userGet.Biography = u.Biography

	if u.BirthDate != nil {
		age := s.app.GetAge(r.Context(), *u.BirthDate)
		userGet.Age = &age
		birthdate := u.BirthDate.Format(DateLayout)
		userGet.Birthdate = &birthdate
	}
	s.writeHTTPJsonOK(w, userGet)
}

func (s *Server) UserSearchGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Social Network")
}

func (s *Server) DialogUserIdListGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) DialogUserIdSendPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) FriendDeleteUserIdPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) FriendSetUserIdPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostCreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostDeleteIdPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostFeedGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostGetIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostUpdatePut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
