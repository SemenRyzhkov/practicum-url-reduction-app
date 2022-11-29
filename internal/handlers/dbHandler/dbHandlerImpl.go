package dbHandler

import (
	"database/sql"
	"net/http"
)

type dbHandlerImpl struct {
	address string
}

func (d dbHandlerImpl) PingConnection(writer http.ResponseWriter, request *http.Request) {
	db, connectionErr := sql.Open("postgres", d.address)
	if connectionErr != nil {
		http.Error(writer, connectionErr.Error(), http.StatusBadRequest)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		http.Error(writer, pingErr.Error(), http.StatusBadRequest)
	}
	writer.WriteHeader(http.StatusOK)
}

func NewDBHandler(address string) DBHandler {
	return &dbHandlerImpl{address: address}
}
