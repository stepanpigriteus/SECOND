package router

import (
	"fmt"
	"net/http"

	"1337b04rd/internal/adapter/http/handler"
	"1337b04rd/internal/adapter/http/middleware"
	"1337b04rd/internal/domain/service"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []func(http.Handler) http.Handler
}

func RegisterRoutes(mux *http.ServeMux, handlers *handler.AllHandlers, sessionService service.SessionService, userService service.UserService) {
	fmt.Println(">>> Registering route: POST /posts")
	routes := []Route{
		// Роут для поста
		{Method: http.MethodPost, Path: "/post", Handler: handlers.Post.CreatePost, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodGet, Path: "/post/", Handler: handlers.Post.GetPostByID, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodPut, Path: "/post/update", Handler: handlers.Post.UpdatePostPostByID, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodDelete, Path: "/post/delete", Handler: handlers.Post.DeletePostPostPostByID, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},

		// Роут для юзера

		{Method: http.MethodPost, Path: "/user", Handler: handlers.User.CreateUser, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodGet, Path: "/user/", Handler: handlers.User.GetUserByID, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodDelete, Path: "/user/delete", Handler: handlers.User.DeleteUser, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},

		// Роут для комментов
		{Method: http.MethodPost, Path: "/post/comment", Handler: handlers.Comment.CreateComment, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodDelete, Path: "/post/comment/delete", Handler: handlers.Comment.DeleteComment, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},

		// Лядские темплейты
		{Method: http.MethodGet, Path: "/catalog", Handler: handlers.Post.GetPostsHandler, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodGet, Path: "/create-post", Handler: handlers.Post.CreatePostPage, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},

		// Роут для хреновых сессий
		{Method: http.MethodDelete, Path: "/sessions/delete/", Handler: handlers.Session.DeleteSession, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodGet, Path: "/sessions/", Handler: handlers.Session.GetSessionByID, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
		{Method: http.MethodPost, Path: "/sessions", Handler: handlers.Session.CreateSession, Middlewares: []func(http.Handler) http.Handler{middleware.LoggerMiddleware, middleware.SessionMiddleware(sessionService, userService)}},
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
