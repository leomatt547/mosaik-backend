package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/auth"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"

	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	parent := models.Parent{}
	child := models.Child{}

	err = json.Unmarshal(body, &parent)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = json.Unmarshal(body, &child)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	parent.Prepare()
	err = parent.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	child.Prepare()
	err = child.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.ParentSignIn(parent.Email, parent.Password)
	if err != nil {
		token, err := server.ChildSignIn(child.Email, child.Password)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
			return
		}
		responses.JSON(w, http.StatusOK, token)
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) ParentSignIn(email, password string) (string, error) {
	var err error

	parent := models.Parent{}

	err = server.DB.Debug().Model(models.Parent{}).Where("email = ?", email).Take(&parent).Error
	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(parent.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateTokenParent(uint32(parent.ID))
}

func (server *Server) ChildSignIn(email, password string) (string, error) {
	var err error

	child := models.Child{}

	err = server.DB.Debug().Model(models.Child{}).Where("email = ?", email).Take(&child).Error
	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(child.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateTokenChild(uint32(child.ID))
}
