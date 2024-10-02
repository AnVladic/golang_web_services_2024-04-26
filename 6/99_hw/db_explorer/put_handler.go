package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (h *HandlerExplorer) createRecord(
	w http.ResponseWriter, r *http.Request) {

	trimmedPath := strings.Trim(r.URL.Path, "/")
	if trimmedPath == "" {
		h.getTablesList(w, r)
		return
	}

	parts := strings.Split(trimmedPath, "/")
	if len(parts) != 1 {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: "unknown path"})
		return
	}
	tableName := parts[0]
	tableMeta := h.tables[tableName]
	if tableMeta == nil {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: "unknown table"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var t map[string]interface{}
	err := decoder.Decode(&t)
	if err != nil {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: err.Error()})
		return
	}

	var columns []string
	var placeholders []string
	var values []interface{}
	for col, val := range t {
		if col == tableMeta.primaryKey {
			continue
		}
		field := tableMeta.fields[col]
		if field.name == "" {
			continue
		}

		columns = append(columns, col)
		placeholders = append(placeholders, "?")
		values = append(values, val)
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	result, err := h.db.Exec(query, values...)
	fmt.Println(columns, values, err)
	if err != nil {
		MarshalAndWrite(w, http.StatusInternalServerError, ResultResponse{Error: err.Error()})
		return
	}
	lastInsertID, _ := result.LastInsertId()
	MarshalAndWrite(w, http.StatusOK, ResultResponse{
		Response: map[string]int64{tableMeta.primaryKey: lastInsertID},
	})
}
