package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/filatkinen/socialnet/internal/common"
	socialapp "github.com/filatkinen/socialnet/internal/server/app"
	"github.com/filatkinen/socialnet/internal/storage"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

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
		t, err := time.Parse(common.DateLayout, *userReg.Birthdate)
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
		birthdate := u.BirthDate.Format(common.DateLayout)
		userGet.Birthdate = &birthdate
	}
	s.writeHTTPJsonOK(w, userGet)
}

func (s *Server) UserSearchGet(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	first_name := r.URL.Query().Get("first_name")
	second_name := r.URL.Query().Get("second_name")
	if len(first_name) == 0 && len(second_name) == 0 {
		s.ClientError(w, http.StatusBadRequest, "FirstName and SecondName masks are empty")
		return
	}
	usersGet, err := s.app.UserSearch(r.Context(), first_name, second_name)
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

	var users []*User
	for i := range usersGet {
		var user User
		user.Biography = usersGet[i].Biography
		user.FirstName = usersGet[i].FirstName
		user.SecondName = usersGet[i].SecondName
		if usersGet[i].BirthDate != nil {
			age := s.app.GetAge(r.Context(), *usersGet[i].BirthDate)
			user.Age = &age
			birthdate := usersGet[i].BirthDate.Format(common.DateLayout)
			user.Birthdate = &birthdate
		}
		user.City = usersGet[i].City
		users = append(users, &user)
	}
	s.writeHTTPJsonOK(w, users)
}

func (s *Server) PostFeedGet(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	userID := r.Context().Value(ContextUserKey).(string)
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit == 0 {
		limit = 10
	}
	var postsGet []*storage.Post
	var err error
	if s.cache.Ready() {
		postsGet = s.cache.UserGetFriendsPosts(userID, offset, limit)
	} else {
		postsGet, err = s.app.UserGetFriendsPosts(r.Context(), userID, offset, limit)
	}
	if err != nil {
		if errors.Is(err, socialapp.ErrorPostsNotFound) {
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

	var posts []*Post
	for i := range postsGet {
		var post Post
		post.Id = postsGet[i].PostId
		post.Text = postsGet[i].PostText
		post.AuthorUserId = postsGet[i].UserId
		posts = append(posts, &post)
	}
	s.writeHTTPJsonOK(w, posts)
}

func (s *Server) PostCreatePost(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	userID := r.Context().Value(ContextUserKey).(string)

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
	id, err := s.app.UserAddPost(r.Context(), userID, text)
	if err != nil {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: rID,
			Code:      500,
		})
		return
	}
	go func() {
		_ = s.cache.AddPost(&storage.Post{
			PostId:   id,
			UserId:   userID,
			PostText: text,
			PostDate: time.Now(),
		})

	}()

	err = s.PutPostMessageToRabbit(id, text, userID)
	if err != nil {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: rID,
			Code:      500,
		})
		return
	}

	s.writeHTTPJsonOK(w, PostCreateResponse{Id: id})
}

func (s *Server) FriendSetUserIdPut(w http.ResponseWriter, r *http.Request) {
	rID := r.Context().Value(RequestID).(string)
	userID := r.Context().Value(ContextUserKey).(string)
	vars := mux.Vars(r)
	friendID, ok := vars["user_id"]
	if !ok {
		s.ClientError(w, http.StatusBadRequest, "friendID was not set in URL")
		return
	}
	err := s.app.UserAddFriend(r.Context(), userID, friendID)
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

func (s *Server) DialogUserIdListGet(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) DialogUserIdSendPost(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) PostUpdateCache(w http.ResponseWriter, r *http.Request) {
	s.cache.UpdatePostAll()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) ShardUpdate(w http.ResponseWriter, r *http.Request) {
	err := s.app.Storage.GetShards(r.Context())
	if err != nil {
		s.ServerError(w, http.StatusInternalServerError, &InlineResponse500{
			Message:   err.Error(),
			RequestId: "",
			Code:      500,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostFeedPosted(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ContextUserKey).(string)
	err := s.ws.NewConnection(w, r, userID)
	if err != nil {
		s.log.Println("Error creating ws comnection " + err.Error())
	}
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Social Network")
}

func (s *Server) FriendDeleteUserIdPut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Server) PostDeleteIdPut(w http.ResponseWriter, r *http.Request) {
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
