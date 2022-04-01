package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strings"

	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/models"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/responses"
	"gitlab.informatika.org/if3250_2022_37_mosaik/mosaik-backend/src/api/utils/formaterror"
)

func (server *Server) SendMail(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	parent := models.Parent{}

	err = json.Unmarshal(body, &parent)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	parent.Prepare()
	err = parent.Validate("reset")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.DB.Debug().Model(models.Parent{}).Where("email = ?", parent.Email).Take(&parent).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Sender data.
	from := "mosaik.id.noreply@gmail.com"
	password := "mosaik-id-admin"

	// Receiver email address.
	to := []string{parent.Email}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	parentdetail, pw, err := parent.ResetParentPassword(server.DB, uint32(parent.ID))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	// Message.
	var msg strings.Builder
	msg.WriteString("Hi, ")
	msg.WriteString(parent.Nama)
	msg.WriteString(".\n")
	msg.WriteString("We have received your request to reset your password.")
	msg.WriteString("This is your new password: ")
	msg.WriteString(pw)
	msg.WriteString("\n\nBecause you have requested to reset your password, ")
	msg.WriteString("you have to login with the given password. ")
	msg.WriteString("You will be redirected to another page to set your new password after you login.")
	message := []byte(msg.String())

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, parentdetail)
}
