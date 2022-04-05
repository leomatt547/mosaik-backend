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

func (server *Server) SendMailParent(w http.ResponseWriter, r *http.Request) {

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
	msg.WriteString("We have received your request to reset your password.\n")
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

func (server *Server) SendMailChild(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	child := models.Child{}
	parent := models.Parent{}

	err = json.Unmarshal(body, &child)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	child.Prepare()
	err = child.Validate("reset")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.DB.Debug().Model(models.Child{}).Where("email = ?", child.Email).Take(&child).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = server.DB.Debug().Model(models.Parent{}).Where("id = ?", child.ParentID).Take(&child.Parent).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	parent.Prepare()
	err = server.DB.Debug().Model(models.Parent{}).Where("id = ?", child.ParentID).Take(&parent).Error
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Sender data.
	from := "mosaik.id.noreply@gmail.com"
	password := "mosaik-id-admin"

	// Receiver email address.
	to := []string{child.Email}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	childdetail, pw, err := child.ResetChildPassword(server.DB, uint64(parent.ID))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	childdetail.Parent = child.Parent

	// Message.
	var msg1 strings.Builder
	msg1.WriteString("Hi, ")
	msg1.WriteString(child.Nama)
	msg1.WriteString(".\n")
	msg1.WriteString("We have received your request to reset your password.\n")
	msg1.WriteString("This is your new password: ")
	msg1.WriteString(pw)
	msg1.WriteString("\n\nBecause you have requested to reset your password, ")
	msg1.WriteString("you have to ask your parent to change your password from the given password to the new one. ")
	msg1.WriteString("You will not be able to login to your account until your parent change your password. ")
	msg1.WriteString("If your parent have changed your password to the new one, you will be able to login to your account again.")
	message := []byte(msg1.String())

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	// Receiver email address.
	to = []string{parent.Email}

	// Message.
	var msg2 strings.Builder
	msg2.WriteString("Hi, ")
	msg2.WriteString(parent.Nama)
	msg2.WriteString(".\n")
	msg2.WriteString("We want to inform you that your child, ")
	msg2.WriteString(child.Nama)
	msg2.WriteString("(")
	msg2.WriteString(child.Email)
	msg2.WriteString("), have requested to reset their password.\n")
	msg2.WriteString("This is your child new password: ")
	msg2.WriteString(pw)
	msg2.WriteString("\n\nBecause your child have requested to reset their password, ")
	msg2.WriteString("you are required to help them to change their password.")
	msg2.WriteString("Your child will not be able to login to their account until you change it into the new one.\n\n")
	msg2.WriteString("How to change your child password :\n")
	msg2.WriteString("1. Login to your account\n")
	msg2.WriteString("2. Go to \"Manage Child Account\"\n")
	msg2.WriteString("3. Choose your child account\n")
	msg2.WriteString("4. Click the \"Change Password\" button\n")
	msg2.WriteString("5. Fill the \"Old Password\" with the given password above\n")
	msg2.WriteString("6. Fill the \"New Password\" and \"Confirm New Password\" with the new one\n")
	msg2.WriteString("7. Click \"Save\"\n")
	msg2.WriteString("After you changed your child password, your child will be able to login to their account again.")
	message = []byte(msg2.String())

	// Authentication.
	auth = smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}

	responses.JSON(w, http.StatusOK, childdetail)
}
