package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/auth"
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
	vars := r.URL.Query().Get("parent_id")
	pid, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		bl := models.Blacklist{}
		bls, err := bl.FindAllBlacklist(server.DB)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, bls)
	} else {
		// vars2 := r.URL.Query().Get("url")
		// url, err := strconv.Parse
		bl := models.Blacklist{}
		bls, err := bl.FindBlacklistByParentID(server.DB, uint32(pid))
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, bls)
	}
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

func (server *Server) DeleteBlacklist(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	vars := mux.Vars(r)

	// Is a valid blaclist id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenParentID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the blacklist exist
	bl := models.Blacklist{}
	err = server.DB.Debug().Model(models.Blacklist{}).Where("id = ?", pid).Take(&bl).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this child?
	if uid != bl.ParentID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = bl.DeleteBlacklist(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
