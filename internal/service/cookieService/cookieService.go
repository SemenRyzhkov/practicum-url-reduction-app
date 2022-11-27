package cookieService

import "net/http"

type CookieService interface {
	GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(w http.ResponseWriter, r *http.Request, name string) (string, error)
}
