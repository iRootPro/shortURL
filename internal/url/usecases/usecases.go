package usecases

import (
	"encoding/base64"
)

func GenerateShortLink(originalLink []byte) string {
	return base64.StdEncoding.EncodeToString(originalLink)
}
