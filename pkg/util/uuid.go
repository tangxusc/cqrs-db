package util

import (
	"github.com/gofrs/uuid"
)

func GenerateUuid() string {
	u2, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return u2.String()
}
