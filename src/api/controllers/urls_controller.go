package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"

	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
)

func (server *Server) CreateUrl(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	url := models.Url{}
	err = json.Unmarshal(body, &url)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = url.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	urlCreated, err := url.SaveUrl(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, urlCreated.ID))
	responses.JSON(w, http.StatusCreated, urlCreated)
}

func (server *Server) GetUrls(w http.ResponseWriter, r *http.Request) {
	url := models.Url{}
	urls, err := url.FindAllUrls(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, urls)
}

func (server *Server) GetUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	url := models.Url{}
	urlGotten, err := url.FindUrlByID(server.DB, uint64(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, urlGotten)
}
