package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {

	if strings.Contains(err, "fullname") {
		return errors.New("Fullname Already Taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Password")
	}
	return errors.New("incorrect details")
}
