package handler

import (
	"1337b04rd/internal/domain/entity"
	"1337b04rd/internal/domain/service"
	externalfunc "1337b04rd/pkg/external_func"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> CreateUser handler called")
	var user entity.User
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error":"invalid input!"}`, http.StatusBadRequest)
		return
	}

	createdUser, err := h.userService.CreateUser(ctx, user)
	if err != nil {
		http.Error(w, `{"error":"failed to create user"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(createdUser); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> GetUserByID handler called")
	if h.userService == nil {
		http.Error(w, `{"error":"userService is not initialized"}`, http.StatusInternalServerError)
		return
	}
	id, err := externalfunc.ExtractIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	post, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		http.Error(w, `{"error":"user not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> UpdateUser handler called")
	var user entity.User
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error":"invalid input!"}`, http.StatusBadRequest)
		return
	}

	updatedUser, err := h.userService.UpdateUser(ctx, user)
	if err != nil {
		http.Error(w, `{"error":"failed to update user"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(">>> DeleteUser handler called")
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	id, err := externalfunc.ExtractIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}
	if err := h.userService.DeleteUser(ctx, id); err != nil {
		http.Error(w, `{"error":"failed to delete user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
