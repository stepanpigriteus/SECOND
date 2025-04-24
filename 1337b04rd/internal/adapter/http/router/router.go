package router

import (
	"1337b04rd/internal/adapter/http/handler"
	"1337b04rd/internal/adapter/http/middleware"
	"fmt"
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
}

func RegisterRoutes(mux *http.ServeMux, handlers *handler.AllHandlers) {
	fmt.Println(">>> Registering route: POST /posts")
	routes := []Route{
		// Роут для поста
		{Method: http.MethodPost, Path: "/post", Handler: handlers.Post.CreatePost, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodGet, Path: "/post/", Handler: handlers.Post.GetPostByID, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodPut, Path: "/post/update", Handler: handlers.Post.UpdatePostPostByID, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodDelete, Path: "/post/delete", Handler: handlers.Post.DeletePostPostPostByID, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		// Роут для юзера

		{Method: http.MethodPost, Path: "/user", Handler: handlers.User.CreateUser},
		{Method: http.MethodGet, Path: "/user/", Handler: handlers.User.GetUserByID},
		{Method: http.MethodDelete, Path: "/user/delete", Handler: handlers.User.DeleteUser},
	}

	for _, route := range routes {
		finalHandler := applyMiddleware(route.Handler, route.Middlewares)
		mux.Handle(route.Path, methodHandler(route.Method, finalHandler))
	}
}

// ограничивает хендлер только одним методом
func methodHandler(method string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf(">>> methodHandler: got %s request on %s\n", r.Method, r.URL.Path)
		if r.Method != method {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func applyMiddleware(h http.Handler, middlewares []func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
