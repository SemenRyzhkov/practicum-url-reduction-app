package dbHandler

import "net/http"

type DBHandler interface {
	PingConnection(writer http.ResponseWriter, request *http.Request)
}
