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
	_ "github.com/heroku/x/hmetrics/onload"
)

func (server *Server) CreateParentDownload(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	//POST variabel butuh ParentID; URL; Title
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	parentdownload := models.ParentDownload{}
	err = json.Unmarshal(body, &parentdownload)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = parentdownload.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenParentID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != parentdownload.ParentID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	parentCreated, err := parentdownload.SaveParentDownload(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, parentCreated.ID))
	responses.JSON(w, http.StatusCreated, parentCreated)
}

func (server *Server) GetParentDownloads(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	vars := r.URL.Query().Get("parent_id")
	pid, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		//minta seluruhnya
		parentdownload := models.ParentDownload{}

		parents, err := parentdownload.FindAllParentDownloads(server.DB)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, parents)
	} else {
		//query parent_id diterima
		parentdownload := models.ParentDownload{}

		parents, err := parentdownload.FindParentDownloadsbyParentID(server.DB, uint32(pid))
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, parents)
	}
}

func (server *Server) GetParentDownload(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	parentdownload := models.ParentDownload{}

	parentReceived, err := parentdownload.FindParentDownloadByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, parentReceived)
}

func (server *Server) DeleteParentDownload(w http.ResponseWriter, r *http.Request) {
	//PERHATIAN: Hanya Parent yang Bisa Hapus!
	//cors.EnableCors(&w)
	vars := mux.Vars(r)

	// Is a valid parent download id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this parent authenticated by Parent?
	_, err = auth.ExtractTokenParentID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the parent exist
	parentdownload := models.ParentDownload{}
	err = server.DB.Debug().Model(models.ParentDownload{}).Where("id = ?", pid).Take(&parentdownload).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}
	_, err = parentdownload.DeleteAParentDownload(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
