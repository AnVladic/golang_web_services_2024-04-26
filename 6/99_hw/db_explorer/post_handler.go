package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (h *HandlerExplorer) changeRecord(
	w http.ResponseWriter, r *http.Request) {

	trimmedPath := strings.Trim(r.URL.Path, "/")
	if trimmedPath == "" {
		h.getTablesList(w, r)
		return
	}

	parts := strings.Split(trimmedPath, "/")
	if len(parts) != 2 {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: "unknown path"})
	}
	rowId, err := strconv.Atoi(parts[1])
	if err != nil {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: err.Error()})
		return
	}

	tableName := parts[0]
	tableMeta := h.tables[tableName]
	if tableMeta == nil {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: "unknown table"})
		return
	}

	exists, err := recordExists(h.db, tableName, tableMeta.primaryKey, rowId)
	if err != nil {
		MarshalAndWrite(w, http.StatusInternalServerError, ResultResponse{Error: err.Error()})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var t map[string]interface{}
	var setClauses []string
	var args []interface{}

	err = decoder.Decode(&t)
	if err != nil {
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: err.Error()})
		return
	}

	for column, value := range t {
		if exists && column == tableMeta.primaryKey {
			MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{
				Error: fmt.Sprintf("field %s have invalid type", column)})
			return
		}

		field := tableMeta.fields[column]
		if field.name == "" {
			MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: "unknown field"})
			return
		}

		if value == nil && field.required {
			MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{
				Error: fmt.Sprintf("field %s have invalid type", field.name)})
			return
		}

		switch value.(type) {
		case float64:
			if !strings.Contains(field.fieldType, "int") {
				MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{
					Error: fmt.Sprintf("field %s have invalid type", field.name)})
				return
			}
		default:
			if strings.Contains(field.fieldType, "int") {
				MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{
					Error: fmt.Sprintf("field %s have invalid type", field.name)})
				return
			}
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
		args = append(args, value)
	}

	// Добавляем условие WHERE
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s = ?",
		tableName, strings.Join(setClauses, ", "), tableMeta.primaryKey)
	args = append(args, rowId)
	_, err = h.db.Exec(query, args...)
	if err != nil {
		fmt.Println(query, args, err.Error())
		MarshalAndWrite(w, http.StatusBadRequest, ResultResponse{Error: err.Error()})
		return
	}
	MarshalAndWrite(w, http.StatusOK, ResultResponse{Response: map[string]int{"updated": 1}})
}
