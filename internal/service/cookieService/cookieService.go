package cookieService

import "net/http"

type CookieService interface {
	WriteSigned(w http.ResponseWriter) error
	ReadSigned(r *http.Request, name string) (string, error)
}
