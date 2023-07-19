package internalhttp

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (s *Server) NewRouter() http.Handler {
	type Route struct {
		Name          string
		Method        string
		Pattern       string
		HandlerFunc   http.HandlerFunc
		NeedAuthorize bool
	}

	type Routes []Route

	var routes = Routes{
		Route{
			"Index",
			"GET",
			"/",
			s.Index,
			false,
		},

		Route{
			"DialogUserIdListGet",
			strings.ToUpper("Get"),
			"/dialog/{user_id}/list",
			s.DialogUserIdListGet,
			true,
		},

		Route{
			"DialogUserIdSendPost",
			strings.ToUpper("Post"),
			"/dialog/{user_id}/send",
			s.DialogUserIdSendPost,
			true,
		},

		Route{
			"FriendDeleteUserIdPut",
			strings.ToUpper("Put"),
			"/friend/delete/{user_id}",
			s.FriendDeleteUserIdPut,
			true,
		},

		Route{
			"FriendSetUserIdPut",
			strings.ToUpper("Put"),
			"/friend/set/{user_id}",
			s.FriendSetUserIdPut,
			true,
		},

		Route{
			"LoginPost",
			strings.ToUpper("Post"),
			"/login",
			s.LoginPost,
			false,
		},

		Route{
			"PostCreatePost",
			strings.ToUpper("Post"),
			"/post/create",
			s.PostCreatePost,
			true,
		},

		Route{
			"PostDeleteIdPut",
			strings.ToUpper("Put"),
			"/post/delete/{id}",
			s.PostDeleteIdPut,
			true,
		},

		Route{
			"PostFeedGet",
			strings.ToUpper("Get"),
			"/post/feed",
			s.PostFeedGet,
			true,
		},

		Route{
			"PostGetIdGet",
			strings.ToUpper("Get"),
			"/post/get/{id}",
			s.PostGetIdGet,
			false,
		},

		Route{
			"PostUpdatePut",
			strings.ToUpper("Put"),
			"/post/update",
			s.PostUpdatePut,
			true,
		},

		Route{
			"UserGetIdGet",
			strings.ToUpper("Get"),
			"/user/get/{id}",
			s.UserGetIdGet,
			false,
		},

		Route{
			"UserRegisterPost",
			strings.ToUpper("Post"),
			"/user/register",
			s.UserRegisterPost,
			false,
		},

		Route{
			"UserSearchGet",
			strings.ToUpper("Get"),
			"/user/search",
			s.UserSearchGet,
			false,
		},
	}

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		if route.NeedAuthorize {
			handler = s.addRequestCounter(s.Logger(s.CheckSession(handler), route.Name))
		} else {
			handler = s.addRequestCounter(s.Logger(handler, route.Name))
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}