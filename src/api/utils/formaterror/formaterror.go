package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "email") {
		return errors.New("email sudah diambil")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("incorrect password")
	}
	if strings.Contains(err, "record not found") {
		return errors.New("incorrect details")
	}
	return errors.New("incorrect details")
}
