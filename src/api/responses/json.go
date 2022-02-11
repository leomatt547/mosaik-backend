package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/cors"
)

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	//cors.EnableCors(&w)
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

func ERROR(w http.ResponseWriter, statusCode int, err error) {
	//cors.EnableCors(&w)
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}