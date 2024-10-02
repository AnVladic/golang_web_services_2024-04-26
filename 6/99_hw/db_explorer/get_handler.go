package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (h *HandlerExplorer) getTablesList(w http.ResponseWriter, r *http.Request) {
	tables := make([]string, len(h.tables))
	i := 0
	for table := range h.tables {
		tables[i] = table
		i++
	}

	MarshalAndWrite(w, http.StatusOK, ResultResponse{
		Response: map[string][]string{"tables": tables}})
}

func (h *HandlerExplorer) getTableRecordsList(
	w http.ResponseWriter, r *http.Request, tableName string) {
	if !contains(h.tables, tableName) {
		MarshalAndWrite(w, http.StatusNotFound, ResultResponse{Error: "unknown table"})
		return
	}

	queryRequest := r.URL.Query()
	limit, err := tryGetInt(&queryRequest, "limit", 5)
	if err != nil {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: err.Error()})
		return
	}
	offset, err := tryGetInt(&queryRequest, "offset", 0)
	if err != nil {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: err.Error()})
		return
	}

	query := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", tableName)
	rows, err := h.db.Query(query, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	results := getRows(w, rows)
	MarshalAndWrite(w, http.StatusOK, ResultResponse{Response: map[string]interface{}{"records": results}})
}

func (h *HandlerExplorer) getTableRow(
	w http.ResponseWriter, r *http.Request, tableName string, rowId int) {

	tableMeta := h.tables[tableName]
	if tableMeta == nil {
		MarshalAndWrite(w, http.StatusNotFound, ResultResponse{Error: "unknown table"})
		return
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, tableMeta.primaryKey)
	rows, err := h.db.Query(query, rowId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	results := getRows(w, rows)
	if len(results) >= 1 {
		MarshalAndWrite(
			w, http.StatusOK, ResultResponse{Response: map[string]interface{}{"record": results[0]}})
		return
	}
	MarshalAndWrite(w, http.StatusNotFound, ResultResponse{Error: "record not found"})
}

func (h *HandlerExplorer) getRequest(w http.ResponseWriter, r *http.Request) {
	trimmedPath := strings.Trim(r.URL.Path, "/")
	if trimmedPath == "" {
		h.getTablesList(w, r)
		return
	}

	parts := strings.Split(trimmedPath, "/")

	switch len(parts) {
	case 1:
		h.getTableRecordsList(w, r, parts[0])
	case 2:
		rowId, err := strconv.Atoi(parts[1])
		if err != nil {
			MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: err.Error()})
		}
		h.getTableRow(w, r, parts[0], rowId)
	default:
		h.getTablesList(w, r)
	}
}
