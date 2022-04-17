package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/auth"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"
)

type WhitelistResult struct {
	URL       string `json:"url"`
	IsAllowed bool   `json:"is_allowed"`
}

func (server *Server) CreateWhitelist(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	wl := models.Whitelist{}
	err = json.Unmarshal(body, &wl)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)

	domain := re.FindAllString(wl.Url, -1)
	for _, element := range domain {
		wl.Url = string(element)
	}

	err = wl.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenParentID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	wl.ParentID = uid

	wlCreated, err := wl.SaveWhitelist(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, wlCreated.ID))
	responses.JSON(w, http.StatusCreated, wlCreated)
}

func (server *Server) GetWhitelist(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query().Get("parent_id")
	pid, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		wl := models.Whitelist{}
		wls, err := wl.FindAllWhitelist(server.DB)
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, wls)
	} else {
		wl := models.Whitelist{}
		wls, err := wl.FindWhitelistByParentID(server.DB, uint32(pid))
		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}
		responses.JSON(w, http.StatusOK, wls)
	}
}

func (server *Server) GetWhitelistByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	wl := models.Whitelist{}
	wls, err := wl.FindWhitelistByID(server.DB, uint64(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, wls)
}

func (server *Server) DeleteWhitelist(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	vars := mux.Vars(r)

	// Is a valid wlaclist id given to us?
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

	// Check if the Whitelist exist
	wl := models.Whitelist{}
	err = server.DB.Debug().Model(models.Whitelist{}).Where("id = ?", pid).Take(&wl).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this child?
	if uid != wl.ParentID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = wl.DeleteWhitelist(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}

func (server *Server) WhitelistChecker(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	wl := models.Whitelist{}

	err = json.Unmarshal(body, &wl)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	hasil_final := WhitelistResult{}
	hasil_final.URL = wl.Url

	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)

	uid, err := auth.ExtractTokenChildID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	domain := re.FindAllString(wl.Url, -1)
	for _, element := range domain {
		_, err = wl.FindRecordByUrl(server.DB, element, uid)
	}
	if err != nil {
		//List Block belum tercantum belum ada
		hasil_final.IsAllowed = false
		responses.JSON(w, http.StatusOK, hasil_final)
		return
	} else {
		//List block sudah ada, perbolehkan
		hasil_final.IsAllowed = true
		responses.JSON(w, http.StatusOK, hasil_final)
		return
	}
}
