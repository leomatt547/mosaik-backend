package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/auth"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"

	"github.com/gorilla/mux"
)

func (server *Server) CreateChildVisit(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	//POST variabel butuh ChildID; URL; Title
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	childvisit := models.ChildVisit{}
	err = json.Unmarshal(body, &childvisit)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = childvisit.Validate()
	if err != nil {
		//dapatkan url id
		url_id, err2 := childvisit.Url.FindRecordByUrl(server.DB, childvisit.Url.Url)
		childvisit.UrlID = url_id.ID
		err = childvisit.Validate()
		if err2 != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
	}
	uid, err := auth.ExtractTokenChildID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != childvisit.ChildID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	childCreated, err := childvisit.SaveChildVisit(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, childCreated.ID))
	responses.JSON(w, http.StatusCreated, childCreated)
}

func (server *Server) GetChildVisits(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	vars := r.URL.Query().Get("child_id")
	cid, err := strconv.ParseUint(vars, 10, 64)
	if err != nil {
		//minta seluruhnya
		childvisit := models.ChildVisit{}

		childs, err := childvisit.FindAllChildVisits(server.DB)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, childs)
	} else {
		//query child_id diterima
		childvisit := models.ChildVisit{}

		childs, err := childvisit.FindChildVisitsbyChildID(server.DB, cid)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, childs)
	}
}

func (server *Server) GetChildVisit(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	childvisit := models.ChildVisit{}

	childReceived, err := childvisit.FindChildVisitByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, childReceived)
}

func (server *Server) DeleteChildVisit(w http.ResponseWriter, r *http.Request) {
	//PERHATIAN: Hanya Parent yang Bisa Hapus!
	//cors.EnableCors(&w)
	vars := mux.Vars(r)

	// Is a valid child visit id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this child authenticated by Parent?
	uid, err := auth.ExtractTokenParentID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the child exist
	childvisit := models.ChildVisit{}
	err = server.DB.Debug().Model(models.ChildVisit{}).Where("id = ?", pid).Take(&childvisit).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	child := models.Child{}
	err = server.DB.Debug().Model(models.Child{}).Where("id = ?", childvisit.ChildID).Take(&child).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	childvisit.Child = child
	// Is the authenticated user, the owner of this child?
	if uid != childvisit.Child.ParentID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = childvisit.DeleteAChildVisit(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
