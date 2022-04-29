package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/notifications"

	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
)

type Result struct {
	URL       string `json:"url"`
	IsBlocked bool   `json:"is_blocked"`
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
	start_time := time.Now()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	type NSFW_temp struct {
		Url       string `json:"url"`
		ChildID   uint64 `json:"child_id"`
		IsBlocked bool   `json:"is_blocked"`
	}

	nsfw := NSFW_temp{}
	model_nsfw := models.NSFWUrl{}

	err = json.Unmarshal(body, &nsfw)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	hasil_final := Result{}
	hasil_final.URL = nsfw.Url

	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)

	domain := re.FindAllStringSubmatch(nsfw.Url, -1)
	for _, element := range domain {
		_, err = model_nsfw.FindRecordByNSFWUrl(server.DB, element[1])
	}
	if err != nil {
		//List Block belum tercantum belum ada, mari memfilter
		// Request the HTML page.
		resp, err := http.Get(nsfw.Url)
		if err != nil {
			hasil_final.IsBlocked = false
			responses.JSON(w, http.StatusOK, hasil_final)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			hasil_final.IsBlocked = false
			responses.JSON(w, http.StatusOK, hasil_final)
			return
		}

		htmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			hasil_final.IsBlocked = false
			responses.JSON(w, http.StatusOK, hasil_final)
			return
		}

		imageRegExp := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)

		subMatchSlice := imageRegExp.FindAllStringSubmatch(string(htmlData), -1)
		for _, item := range subMatchSlice {
			//cek bila timeout
			duration := time.Since(start_time)
			if duration.Seconds() > 28 {
				fmt.Println(duration.Seconds())
				//asumsi konten semuanya baik
				hasil_final.IsBlocked = false
				responses.JSON(w, http.StatusOK, hasil_final)
				return
			}

			//kalimat = append(kalimat, "Image found : "+item[1])
			url := "https://mosaik-ai.herokuapp.com/api/image/classify?url=" + string(item[1])
			method := "GET"
			fmt.Println(url)

			client := &http.Client{
				Timeout: 300 * time.Second,
			}
			req, err := http.NewRequest(method, url, nil)

			if err != nil {
				hasil_final.IsBlocked = false
				responses.JSON(w, http.StatusOK, hasil_final)
				return
			}
			res, err := client.Do(req)
			if err != nil {
				hasil_final.IsBlocked = false
				responses.JSON(w, http.StatusOK, hasil_final)
				return
			}
			defer res.Body.Close()

			res_ai := AIResult{}
			// fmt.Println(res.Body)
			err = json.NewDecoder(res.Body).Decode(&res_ai)
			if err != nil {
				// responses.ERROR(w, http.StatusUnprocessableEntity, err)
				continue
				// return
			}
			for _, hasil := range res_ai {
				if hasil.ClassName == "Porn" || hasil.ClassName == "Sexy" || hasil.ClassName == "Hentai" {
					if hasil.Probability >= 0.5 {
						hasil_final.IsBlocked = true
						saved_url := models.NSFWUrl{}
						//extract domain
						re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)

						submatchall := re.FindAllStringSubmatch(nsfw.Url, -1)
						for _, element := range submatchall {
							saved_url.Url = string(element[1])
						}
						//save hasil ke model NSFW
						err = saved_url.Validate()
						if err != nil {
							hasil_final.IsBlocked = false
							responses.JSON(w, http.StatusOK, hasil_final)
							return
						}
						_, _ = saved_url.SaveNSFWUrl(server.DB)
						child := models.Child{}
						childReceived, err := child.FindChildByID(server.DB, nsfw.ChildID)
						if err != nil {
							responses.ERROR(w, http.StatusInternalServerError, err)
							return
						}
						parentReceived := childReceived.Parent
						fcm_token := parentReceived.FCM
						notifications.SendPushNotification("Akses Konten Ilegal",
							"Anak anda yang bernama "+childReceived.Nama+" telah membuka konten negatif pada link "+nsfw.Url+" dengan waktu "+time.Now().String(), fcm_token)
						responses.JSON(w, http.StatusOK, hasil_final)
						return
					}
				}
			}
		}
		//Tidak ada gambar atau memang web bersih
		hasil_final.IsBlocked = false
		responses.JSON(w, http.StatusOK, hasil_final)
		return
	} else {
		//List block sudah ada, gas block aja
		hasil_final.IsBlocked = true
		child := models.Child{}
		childReceived, err := child.FindChildByID(server.DB, nsfw.ChildID)

		if err != nil {
			responses.ERROR(w, http.StatusInternalServerError, err)
			return
		}

		//Dapatkan token
		parentReceived := childReceived.Parent
		fcm_token := parentReceived.FCM
		notifications.SendPushNotification("Akses Konten Ilegal",
			"Anak anda yang bernama "+childReceived.Nama+" telah membuka konten negatif pada link "+nsfw.Url+" dengan waktu "+time.Now().String(), fcm_token)
		responses.JSON(w, http.StatusOK, hasil_final)
		return
	}
}

func (server *Server) GetImagesInLink(w http.ResponseWriter, r *http.Request) {
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
