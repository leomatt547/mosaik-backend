package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"
)

func (server *Server) CreateBlackList(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	bl := models.Blacklist{}
	err = json.Unmarshal(body, &bl)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = bl.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	blCreated, err := bl.SaveBlacklist(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, blCreated.ID))
	responses.JSON(w, http.StatusCreated, blCreated)
}

func (server *Server) GetBlacklist(w http.ResponseWriter, r *http.Request) {
	bl := models.Blacklist{}
	bls, err := bl.FindAllBlacklist(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, bls)
}

func (server *Server) GetBlacklistByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	bl := models.Blacklist{}
	bls, err := bl.FindBlacklistByID(server.DB, uint64(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, bls)
}
