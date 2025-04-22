package router

import (
	"1337b04rd/internal/adapter/http/handler"
	"net/http"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
}

func RegisterRoutes(mux *http.ServeMux, handlers *handler.AllHandlers) {
	routes := []Route{
		// Роут для создания поста
		{Method: http.MethodPost, Path: "/posts", Handler: handlers.Post.CreatePost},
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
func methodHandler(method string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	}
}

func applyMiddleware(h http.Handler, middlewares []func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
