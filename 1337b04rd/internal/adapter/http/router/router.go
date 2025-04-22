package router

import (
	"fmt"
	"net/http"

	"1337b04rd/internal/adapter/http/handler"
	"1337b04rd/internal/adapter/http/middleware"
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
		// Роут для создания поста
		{Method: http.MethodPost, Path: "/posts", Handler: handlers.Post.CreatePost, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		{Method: http.MethodPost, Path: "/posts/", Handler: handlers.Post.CreatePost, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware}},
		// Пример других роутов
		// {Method: http.MethodGet, Path: "/posts", Handler: handlers.Post.GetPosts},
		// {Method: http.MethodGet, Path: "/posts/{id}", Handler: handlers.Post.GetPostByID},
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
