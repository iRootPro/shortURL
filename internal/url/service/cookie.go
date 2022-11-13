package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var key = []byte("secret key cookie")

func SetCookie() *http.Cookie {
	id := uuid.NewString()

	h := hmac.New(sha256.New, key)
	h.Write([]byte(id))
	dst := h.Sum(nil)

	cookie := &http.Cookie{
		Name:  "token",
		Value: fmt.Sprintf("%s:%x", id, dst),
		Path:  "/",
	}
	return cookie
}

func CheckCookie(cookie *http.Cookie) bool {
	values := strings.Split(cookie.Value, ":")

	data, err := hex.DecodeString(values[1])
	if err != nil {
		log.Fatal(err)
		return false
	}

	id := values[0]

	h := hmac.New(sha256.New, key)
	h.Write([]byte(id))
	sign := h.Sum(nil)

	if hmac.Equal(sign, data) {
		return true
	} else {
		return false
	}
}
