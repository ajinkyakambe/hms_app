package controllers

import (
	"net/http"

	"github.com/repoerna/hms_app/api/handlers"
)

// Home ...
func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	handlers.ResponseJSON(w, http.StatusOK, "HMS APP API - Server Live")
}
