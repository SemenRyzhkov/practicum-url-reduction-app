package cookieservice

import "net/http"

// CookieService интерфейс для куки серввиса
type CookieService interface {
	GetUserIDWithCheckCookieAndIssueNewIfCookieIsMissingOrInvalid(w http.ResponseWriter, r *http.Request, name string) (string, error)
}
