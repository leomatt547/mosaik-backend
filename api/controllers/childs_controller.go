package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/auth"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/utils/cors"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/utils/formaterror"

	"github.com/gorilla/mux"
)

func (server *Server) CreateChild(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	child := models.Child{}
	err = json.Unmarshal(body, &child)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	child.Prepare()
	err = child.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != child.ParentID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	childCreated, err := child.SaveChild(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, childCreated.ID))
	responses.JSON(w, http.StatusCreated, childCreated)
}

func (server *Server) GetChilds(w http.ResponseWriter, r *http.Request) {

	child := models.Child{}

	childs, err := child.FindAllChilds(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, childs)
}

func (server *Server) GetChild(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	child := models.Child{}

	childReceived, err := child.FindChildByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, childReceived)
}

func (server *Server) UpdateChild(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w)
	vars := mux.Vars(r)

	// Check if the child id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the child exist
	child := models.Child{}
	err = server.DB.Debug().Model(models.Child{}).Where("id = ?", pid).Take(&child).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("child not found"))
		return
	}

	// If a user attempt to update a child not belonging to him
	if uid != child.ParentID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the child data
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	childUpdate := models.Child{}
	err = json.Unmarshal(body, &childUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != childUpdate.ParentID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	childUpdate.Prepare()
	err = childUpdate.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	childUpdate.ID = child.ID //this is important to tell the model the child id to update, the other update field are set above

	childUpdated, err := childUpdate.UpdateAChild(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, childUpdated)
}

func (server *Server) DeleteChild(w http.ResponseWriter, r *http.Request) {
	cors.EnableCors(&w)
	vars := mux.Vars(r)

	// Is a valid child id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the child exist
	child := models.Child{}
	err = server.DB.Debug().Model(models.Child{}).Where("id = ?", pid).Take(&child).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this child?
	if uid != child.ParentID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = child.DeleteAChild(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}