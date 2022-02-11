package middlewares

import (
	"errors"
	"net/http"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/auth"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/api/utils/cors"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cors.EnableCors(&w)
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cors.EnableCors(&w)
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}