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

type Result struct {
	URL       string
	isBlocked bool
}

type AIResult []struct {
	ClassName   string  `json:"className"`
	Probability float64 `json:"probability"`
}

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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	nsfw := models.NSFWUrl{}
	err = json.Unmarshal(body, &nsfw)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

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

	hasil_final := Result{}
	hasil_final.URL = nsfw.Url
	subMatchSlice := imageRegExp.FindAllStringSubmatch(string(htmlData), -1)
	for _, item := range subMatchSlice {
		//kalimat = append(kalimat, "Image found : "+item[1])
		url := "https://mosaik-ai.herokuapp.com/api/image/classify?url=" + string(item[1])
		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		defer res.Body.Close()

		res_ai := AIResult{}
		err = json.NewDecoder(res.Body).Decode(&res_ai)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		// fmt.Printf("%s", j)
		for _, hasil := range res_ai {
			if hasil.ClassName == "Porn" || hasil.ClassName == "Sexy" || hasil.ClassName == "Hentai" {
				if hasil.Probability >= 30 {
					hasil_final.isBlocked = true
					responses.JSON(w, http.StatusOK, hasil_final)
					return
				}
			}
		}
		// log.Println("Image found : ", item[1])
	}
	hasil_final.isBlocked = false
	responses.JSON(w, http.StatusOK, hasil_final)
}
