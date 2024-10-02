package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func recordExists(db *sql.DB, table string, primaryKey string, id int) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s = ?)", table, primaryKey)
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func findPrimaryKey(db *sql.DB, dbName, tableName string) (string, error) {
	query := `
		SELECT COLUMN_NAME
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND CONSTRAINT_NAME = 'PRIMARY'
	`

	var primaryKey string
	row := db.QueryRow(query, dbName, tableName)
	err := row.Scan(&primaryKey)
	if err != nil {
		return "", err
	}
	return primaryKey, nil
}

func getFieldsList(db *sql.DB, tableName string) ([]Field, error) {
	query := fmt.Sprintf("SHOW COLUMNS FROM %s", tableName)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	fields := make([]Field, 0)
	defer rows.Close()

	fmt.Printf("Columns in table %s:\n", tableName)
	for rows.Next() {
		var field, colType, null, key, defaultValue, extra string
		_ = rows.Scan(&field, &colType, &null, &key, &defaultValue, &extra)
		filedS := Field{
			name:      field,
			fieldType: colType,
		}

		if null == "NO" {
			filedS.required = true
		}

		fields = append(fields, filedS)
		fmt.Printf("Field: %s, Type: %s, Null: %s\n", field, colType, null)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return fields, nil
}

func contains(slice map[string]*TableMeta, item string) bool {
	val := slice[item]
	return val != nil
}

func tryGetInt(values *url.Values, key string, defaultValue int) (int, error) {
	valStr := values.Get(key)
	if valStr == "" {
		return defaultValue, nil
	}
	if val, err := strconv.Atoi(valStr); err == nil {
		return val, nil
	}
	return defaultValue, nil
}

func getRows(w http.ResponseWriter, rows *sql.Rows) []map[string]interface{} {

	columns, err := rows.Columns()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	columnsType, err := rows.ColumnTypes()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	results := make([]map[string]interface{}, 0)
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		row := make(map[string]interface{}, len(columns))
		for i, col := range columns {
			switch v := values[i].(type) {
			case nil:
				row[col] = nil
			case []byte:
				valStr := string(v)
				switch columnsType[i].DatabaseTypeName() {
				case "INT":
					val, _ := strconv.Atoi(valStr)
					row[col] = val
				default:
					row[col] = valStr
				}
			default:
				row[col] = v
			}
		}
		results = append(results, row)
	}
	return results
}

func MarshalAndWrite(w http.ResponseWriter, status int, response interface{}) {
	result, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(status)
	w.Write(result)
}
