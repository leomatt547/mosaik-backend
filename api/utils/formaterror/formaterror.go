package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "nama") {
		return errors.New("nama sudah ada")
	}

	if strings.Contains(err, "email") {
		return errors.New("email sudah dipakai")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect password")
	}
	return errors.New("incorrect details")
}