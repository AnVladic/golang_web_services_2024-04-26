package main

import (
	"database/sql"
	"net/http"
)

type ResultResponse struct {
	Error    string      `json:"error,omitempty"`
	Response interface{} `json:"response,omitempty"`
}

type TableMeta struct {
	name       string
	primaryKey string
	fields     map[string]Field
}

type HandlerExplorer struct {
	db     *sql.DB
	tables map[string]*TableMeta
}

type Field struct {
	name      string
	fieldType string
	required  bool
}

func (h *HandlerExplorer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getRequest(w, r)
	case http.MethodPut:
		h.createRecord(w, r)
	case http.MethodPost:
		h.changeRecord(w, r)
	case http.MethodDelete:
		h.removeRecord(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Excepted GET, PUT, POST, DELETE but got " + r.Method))
	}
}
