package controllers

import (
	"net/http"

	"gitlab.informatika.org/if3250_2022_37_alkademi/mosaik-backend/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To Mosaik API")
}