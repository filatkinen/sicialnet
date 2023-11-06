package internalhttp

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (s *DialogService) NewRouter() http.Handler {
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
			"LoginPost",
			strings.ToUpper("Post"),
			"/login",
			s.LoginPost,
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
	}

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		if route.Name == "PostFeedPosted" {
			handler = s.CheckSession(handler)
		} else {
			if route.NeedAuthorize {
				handler = s.prometheusMiddleware(s.addRequestCounter(s.Logger(s.CheckSession(handler), route.Name)), route.Pattern)
			} else {
				handler = s.prometheusMiddleware(s.addRequestCounter(s.Logger(handler, route.Name)), route.Pattern)
			}
		}
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.Path("/metrics").Handler(promhttp.Handler())
	return router
}
