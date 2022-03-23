package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"

	"github.com/gorilla/mux"
)

func (server *Server) CreateNSFWUrl(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	url := models.NSFWUrl{}
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
	urlCreated, err := url.SaveNSFWUrl(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, urlCreated.ID))
	responses.JSON(w, http.StatusCreated, urlCreated)
}

func (server *Server) GetNSFWUrls(w http.ResponseWriter, r *http.Request) {
	url := models.NSFWUrl{}
	urls, err := url.FindAllNSFWUrls(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, urls)
}

func (server *Server) GetNSFWUrl(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	url := models.NSFWUrl{}
	urlGotten, err := url.FindNSFWUrlByID(server.DB, uint64(uid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, urlGotten)
}

func (server *Server) SavedSearchChecker(w http.ResponseWriter, r *http.Request) {
	//cors.EnableCors(&w)
	// vars := mux.Vars(r)
	// uid, err := strconv.ParseUint(vars["id"], 10, 32)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusBadRequest, err)
	// 	return
	// }
	var kalimat []string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// parent := models.Parent{}
	nsfw := models.NSFWUrl{}
	err = json.Unmarshal(body, &nsfw)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	kalimat = append(kalimat, "Parsing : "+nsfw.Url)
	// log.Println("Parsing : ", nsfw.Url)

	// Request the HTML page.
	resp, err := http.Get(nsfw.Url)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		responses.ERROR(w, resp.StatusCode, err)
		return
	}

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	imageRegExp := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)

	subMatchSlice := imageRegExp.FindAllStringSubmatch(string(htmlData), -1)
	for _, item := range subMatchSlice {
		kalimat = append(kalimat, "Image found : "+item[1])
		// log.Println("Image found : ", item[1])
	}
	responses.JSON(w, http.StatusOK, kalimat)
}
