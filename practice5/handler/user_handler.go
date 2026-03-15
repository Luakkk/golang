package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"practice5/repository"
)

type Handler struct {
	Repo *repository.Repository
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))

	gender := r.URL.Query().Get("gender")
	orderBy := r.URL.Query().Get("order_by")

	if orderBy == "" {
		orderBy = "name"
	}

	data, err := h.Repo.GetPaginatedUsers(page, size, gender, orderBy)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(data)
}
func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {

	user1 := r.URL.Query().Get("user1")
	user2 := r.URL.Query().Get("user2")

	data, err := h.Repo.GetCommonFriends(user1, user2)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(data)
}