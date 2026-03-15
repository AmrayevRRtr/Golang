package handler

import "net/http"

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.users)
	mux.HandleFunc("/users/", h.userByID)
	mux.HandleFunc("/users/common-friends", h.CommonFriends)
}
