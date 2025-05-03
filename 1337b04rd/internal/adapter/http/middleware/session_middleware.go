package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"a1337b04rd/internal/domain/entity"
	"a1337b04rd/internal/domain/service"
	externalfunc "a1337b04rd/pkg/external_func"
)

func SessionMiddleware(sessionService service.SessionService, userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем все доступные куки для отладки
			cookies := r.Cookies()
			fmt.Printf(">>> [DEBUG] SessionMiddleware: Всего получено кук: %d\n", len(cookies))
			for i, c := range cookies {
				fmt.Printf(">>> [DEBUG] Cookie %d: %s=%s\n", i, c.Name, c.Value)
			}

			// Пытаемся получить cookie session_id
			cookie, err := r.Cookie("session_id")

			// Флаг для отслеживания успешной аутентификации
			authenticated := false

			if err == nil && cookie != nil && cookie.Value != "" {
				fmt.Printf(">>> [DEBUG] SessionMiddleware: Найдена кука session_id=%s\n", cookie.Value)

				// Проверяем сессию в БД
				session, err := sessionService.GetSessionByID(r.Context(), cookie.Value)

				if err == nil && session.ID != "" {
					fmt.Printf(">>> [DEBUG] SessionMiddleware: Сессия найдена, UserID=%d\n", session.UserID)

					// Сессия найдена и действительна
					if session.ExpiresAt.After(time.Now()) {
						// Пытаемся получить пользователя
						user, err := userService.GetUserByID(r.Context(), session.UserID)

						if err == nil {
							fmt.Printf(">>> [DEBUG] SessionMiddleware: Пользователь найден, ID=%d\n", user.ID)

							// Помечаем, что аутентификация успешна
							authenticated = true

							// Обновляем время жизни куки при успешной аутентификации
							http.SetCookie(w, &http.Cookie{
								Name:     "session_id",
								Value:    session.ID,
								Expires:  time.Now().Add(7 * 24 * time.Hour),
								HttpOnly: true,
								Path:     "/",
								SameSite: http.SameSiteLaxMode,
							})

							// Успешная аутентификация - продолжаем запрос
							ctx := context.WithValue(r.Context(), "user", &user)
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						} else {
							fmt.Printf(">>> [DEBUG] SessionMiddleware: Ошибка получения пользователя: %v\n", err)
						}
					} else {
						fmt.Printf(">>> [DEBUG] SessionMiddleware: Сессия истекла: %v > %v\n",
							time.Now(), session.ExpiresAt)
					}
				} else {
					fmt.Printf(">>> [DEBUG] SessionMiddleware: Ошибка получения сессии: %v\n", err)
				}

				// Если дошли до этой точки, сессия недействительна - удаляем куку
				if !authenticated {
					fmt.Println(">>> [DEBUG] SessionMiddleware: Удаляем недействительную куку")
					http.SetCookie(w, &http.Cookie{
						Name:     "session_id",
						Value:    "",
						MaxAge:   -1,
						HttpOnly: true,
						Path:     "/",
						SameSite: http.SameSiteLaxMode,
					})
				}
			} else {
				fmt.Printf(">>> [DEBUG] SessionMiddleware: Кука session_id не найдена: %v\n", err)
			}

			// Если дошли до этой точки и не аутентифицированы - создаем нового пользователя и сессию
			if !authenticated {
				fmt.Println(">>> [DEBUG] SessionMiddleware: Создание нового пользователя")

				// 1. Получаем случайного персонажа
				character, err := externalfunc.GetRandomCharacter()
				if err != nil {
					fmt.Printf(">>> [DEBUG] SessionMiddleware: Ошибка получения персонажа: %v\n", err)
					http.Error(w, "Ошибка при получении персонажа", http.StatusInternalServerError)
					return
				}

				// 2. Создаем нового пользователя
				user := entity.User{
					CharacterName: character.Name,
					AvatarURL:     character.Image,
				}

				// 3. Сохраняем пользователя в БД
				createdUser, err := userService.CreateUser(r.Context(), user)
				if err != nil {
					fmt.Printf(">>> [DEBUG] SessionMiddleware: Ошибка создания пользователя: %v\n", err)
					http.Error(w, "Ошибка при создании пользователя", http.StatusInternalServerError)
					return
				}
				fmt.Printf(">>> [DEBUG] SessionMiddleware: Создан пользователь, ID=%d\n", createdUser.ID)

				// 4. Создаем объект сессии
				newSession := entity.Session{
					UserID: createdUser.ID,
				}

				// 5. Сохраняем сессию в БД через сервис
				createdSession, err := sessionService.CreateSession(r.Context(), newSession)
				if err != nil {
					fmt.Printf(">>> [DEBUG] SessionMiddleware: Ошибка создания сессии: %v\n", err)
					http.Error(w, "Ошибка при создании сессии", http.StatusInternalServerError)
					return
				}
				fmt.Printf(">>> [DEBUG] SessionMiddleware: Создана сессия, ID=%s\n", createdSession.ID)

				// 6. Устанавливаем cookie с ID сессии
				cookieExpiration := time.Now().Add(7 * 24 * time.Hour)
				sessionCookie := &http.Cookie{
					Name:     "session_id",
					Value:    createdSession.ID,
					Expires:  cookieExpiration,
					HttpOnly: true,
					Path:     "/",
					SameSite: http.SameSiteLaxMode,
				}
				http.SetCookie(w, sessionCookie)
				fmt.Printf(">>> [DEBUG] SessionMiddleware: Установлена кука session_id=%s, срок: %v\n",
					createdSession.ID, cookieExpiration)

				// Проверяем, что кука установлена правильно
				fmt.Printf(">>> [DEBUG] SessionMiddleware: Cookie Raw: %s\n", sessionCookie.String())

				// 7. Добавляем пользователя в контекст запроса и продолжаем
				ctx := context.WithValue(r.Context(), "user", createdUser)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		})
	}
}
