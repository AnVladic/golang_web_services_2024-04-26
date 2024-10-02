// тут лежит тестовый код
// менять вам может потребоваться только коннект к базе
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// DSN это соединение с базой
	// вы можете изменить этот на тот который вам нужен
	// docker run -p 3306:3306 -v $(PWD):/docker-entrypoint-initdb.d -e MYSQL_ROOT_PASSWORD=1234 -e MYSQL_DATABASE=golang -d mysql
	// DSN = "root@tcp(localhost:3306)/golang2017?charset=utf8"
	DSN = "root:@tcp(localhost:3306)/coursera?charset=utf8"
)

func getEmptyTableMeta(db *sql.DB) []*TableMeta {
	tables := []*TableMeta{}
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			log.Fatal(err)
		}
		tables = append(tables, &TableMeta{name: tableName})
	}
	return tables
}

func NewDbExplorer(db *sql.DB) (http.Handler, error) {
	handler := &HandlerExplorer{db: db, tables: map[string]*TableMeta{}}
	for _, table := range getEmptyTableMeta(db) {
		table.primaryKey, _ = findPrimaryKey(db, "coursera", table.name)
		table.fields = map[string]Field{}
		fields, _ := getFieldsList(db, table.name)
		for _, field := range fields {
			table.fields[field.name] = field
		}

		handler.tables[table.name] = table
	}

	return handler, nil
}

func main() {
	db, err := sql.Open("mysql", DSN)
	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	fmt.Println(findPrimaryKey(db, "coursera", "items"))
	handler, err := NewDbExplorer(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting server at :8082")
	http.ListenAndServe(":8082", handler)
}
