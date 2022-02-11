package controllers

import (
	"net/http"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	responses.JSON(w, http.StatusOK, "Welcome To Mosaik API")
}