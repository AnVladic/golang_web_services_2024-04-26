package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (h *HandlerExplorer) removeRecord(
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

	result, err := h.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName), rowId)
	if err != nil {
		MarshalAndWrite(w, http.StatusInternalServerError, ResultResponse{Error: err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		MarshalAndWrite(w, http.StatusInternalServerError, ResultResponse{Error: err.Error()})
		return
	}

	MarshalAndWrite(w, http.StatusOK, ResultResponse{Response: map[string]int64{"deleted": rowsAffected}})
}
