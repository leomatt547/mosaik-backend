package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/auth"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"

	"golang.org/x/crypto/bcrypt"
)

type ParentResponse struct {
	ID        uint32
	Nama      string
	Email     string
	Password  string
	IsChange  bool
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Token     string
}

type ChildResponse struct {
	ID        uint64
	Nama      string
	Email     string
	Password  string
	Parent    models.Parent
	ParentID  uint32
	LastLogin time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Token     string
}

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

	parentdetail, err := server.ParentSignIn(parent.Email, parent.Password)
	if err != nil {
		//fmt.Println("errornya di:" + err.Error())
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			formattedError := formaterror.FormatError(err.Error())
			responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
			return
		} else {
			childdetail, err := server.ChildSignIn(child.Email, child.Password)
			if err != nil {
				//fmt.Println("errornya child di:" + err.Error())
				if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
					formattedError := formaterror.FormatError(err.Error())
					responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
					return
				} else if err.Error() == "record not found" {
					formattedError := formaterror.FormatError(err.Error())
					responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
					return
				} else {
					formattedError := formaterror.FormatError(err.Error())
					responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
					return
				}
			} else {
				responses.JSON(w, http.StatusOK, childdetail)
			}
		}
	} else {
		responses.JSON(w, http.StatusOK, parentdetail)
	}
}

func (server *Server) ParentSignIn(email, password string) (*ParentResponse, error) {
	var err error
	parent := models.Parent{}
	response := ParentResponse{}

	err = server.DB.Debug().Model(models.Parent{}).Where("email = ?", email).Take(&parent).Error
	if err != nil {
		return &ParentResponse{}, err
	}

	err = models.VerifyPassword(parent.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &ParentResponse{}, err
	}

	getParent, err := parent.FindParentByID(server.DB, uint32(parent.ID))
	if err != nil {
		return &ParentResponse{}, err
	}

	response.ID = getParent.ID
	response.Email = getParent.Email
	response.Nama = getParent.Nama
	response.Password = getParent.Password
	response.IsChange = getParent.IsChange
	response.LastLogin = getParent.LastLogin
	response.CreatedAt = getParent.CreatedAt
	response.UpdatedAt = getParent.UpdatedAt

	token, err := auth.CreateTokenParent(uint32(parent.ID))
	response.Token = token
	return &response, err
}

func (server *Server) ChildSignIn(email, password string) (*ChildResponse, error) {
	var err error
	child := models.Child{}
	response := ChildResponse{}

	err = server.DB.Debug().Model(models.Child{}).Where("email = ?", email).Take(&child).Error
	if err != nil {
		return &ChildResponse{}, err
	}

	err = models.VerifyPasswordChild(child.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &ChildResponse{}, err
	}

	getChild, err := child.FindChildByID(server.DB, uint64(child.ID))
	if err != nil {
		return &ChildResponse{}, err
	}

	response.ID = getChild.ID
	response.Email = getChild.Email
	response.Nama = getChild.Nama
	response.Password = getChild.Password
	response.Parent = getChild.Parent
	response.ParentID = getChild.Parent.ID
	response.LastLogin = getChild.LastLogin
	response.CreatedAt = getChild.CreatedAt
	response.UpdatedAt = getChild.UpdatedAt

	token, err := auth.CreateTokenChild(uint64(child.ID))
	response.Token = token
	return &response, err
}
