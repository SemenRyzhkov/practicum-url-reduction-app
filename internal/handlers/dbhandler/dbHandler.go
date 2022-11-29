package dbhandler

import "net/http"

type DBHandler interface {
	PingConnection(writer http.ResponseWriter, request *http.Request)
}
