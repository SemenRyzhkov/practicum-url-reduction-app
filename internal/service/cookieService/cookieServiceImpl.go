package cookieService

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

var (
	_               CookieService = &cookieServiceImpl{}
	ErrValueTooLong               = errors.New("cookie value too long")
	ErrInvalidValue               = errors.New("invalid cookie value")
)

type cookieServiceImpl struct {
	secretKey []byte
}

func New(key string) CookieService {
	secretKey, err := hex.DecodeString(key)
	if err != nil {
		log.Fatal(err)
	}
	return &cookieServiceImpl{secretKey}
}

func (c cookieServiceImpl) GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(
	w http.ResponseWriter,
	r *http.Request, name string) (string, error) {
	userID, err := c.readSigned(r, name)
	if err == nil {
		return userID, nil
	}

	if errors.Is(err, http.ErrNoCookie) || errors.Is(err, ErrInvalidValue) {
		userID, writeErr := c.writeSigned(w)
		if writeErr == nil {
			return userID, nil
		}
		return "", writeErr
	}
	return "", err
}

func (c cookieServiceImpl) writeSigned(w http.ResponseWriter) (string, error) {
	userID := uuid.New().String()
	cookie := http.Cookie{
		Name:     "userID",
		Value:    userID,
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
	}
	mac := hmac.New(sha256.New, c.secretKey)
	mac.Write([]byte(cookie.Name))
	mac.Write([]byte(cookie.Value))
	signature := mac.Sum(nil)

	cookie.Value = string(signature) + cookie.Value
	fmt.Println("Write userID " + userID)

	return userID, write(w, cookie)

}

func (c cookieServiceImpl) readSigned(r *http.Request, name string) (string, error) {
	// {signature}{original value}
	signedValue, err := read(r, name)
	if err != nil {
		return "", err
	}

	if len(signedValue) < sha256.Size {
		return "", ErrInvalidValue
	}

	signature := signedValue[:sha256.Size]
	value := signedValue[sha256.Size:]
	fmt.Println("read userID " + value)
	mac := hmac.New(sha256.New, c.secretKey)
	mac.Write([]byte(name))
	mac.Write([]byte(value))
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal([]byte(signature), expectedSignature) {
		return "", ErrInvalidValue
	}

	return value, nil
}

func write(w http.ResponseWriter, cookie http.Cookie) error {
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > 4096 {
		return ErrValueTooLong
	}

	http.SetCookie(w, &cookie)

	return nil
}

func read(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}
