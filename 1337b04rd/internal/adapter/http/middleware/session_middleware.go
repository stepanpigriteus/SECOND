package middleware

import (
	"net/http"
	"time"

	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/service"
	externalfunc "1337b04rd/pkg/external_func"
)

func SessionMiddleware(sessionService service.SessionService, userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err == nil {
				session, err := sessionService.GetSessionByID(r.Context(), cookie.Value)
				if err == nil && session.ExpiresAt.After(time.Now()) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Сессия отсутствует или недействительна — создаём нового пользователя
			character, err := externalfunc.GetRandomCharacter()
			if err != nil {
				http.Error(w, "Ошибка при получении персонажа", http.StatusInternalServerError)
				return
			}

			user := entity.User{
				CharacterName: character.Name,
				AvatarURL:     character.Image,
			}

			createdUser, err := userService.CreateUser(r.Context(), user)
			if err != nil {
				http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
				return
			}

			newSession := entity.Session{
				ID:        externalfunc.GenerateSessionID(),
				UserID:    createdUser.ID,
				ExpiresAt: time.Now().Add(24 * time.Hour),
			}

			createdSession, err := sessionService.CreateSession(r.Context(), newSession)
			if err != nil {
				http.Error(w, "Ошибка при создании сессии", http.StatusInternalServerError)
				return
			}

			// Устанавливаем cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    createdSession.ID,
				Expires:  createdSession.ExpiresAt,
				HttpOnly: true,
				Path:     "/",
			})

			r.AddCookie(&http.Cookie{
				Name:     "session_id",
				Value:    createdSession.ID,
				Expires:  createdSession.ExpiresAt,
				HttpOnly: true,
				Path:     "/",
			})

			next.ServeHTTP(w, r)
		})
	}
}
