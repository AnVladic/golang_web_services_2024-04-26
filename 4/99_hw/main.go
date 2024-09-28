package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

var rootPersons Root

type ParamString struct {
	string
}

type Person struct {
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	About     string `xml:"about"`
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	Gender    string `xml:"gender"`
	Name      string `xml:"-"`
}

type Root struct {
	Persons []Person `xml:"row"`
}

func SortBy(persons *[]Person, field string, by int) error {
	if len(field) == 0 {
		field = "Name"
	}

	if by != OrderByDesc && by != OrderByAsc && by != OrderByAsIs {
		return fmt.Errorf("sort by %s failed", field)
	}

	if field == "Id" {
		if by == OrderByAsc {
			sort.Slice(*persons, func(i, j int) bool {
				return (*persons)[i].Id < (*persons)[j].Id
			})
		} else if by == OrderByDesc {
			sort.Slice(*persons, func(i, j int) bool {
				return (*persons)[i].Id > (*persons)[j].Id
			})
		}
	} else if field == "Age" {
		if by == OrderByAsc {
			sort.Slice(*persons, func(i, j int) bool {
				return (*persons)[i].Age < (*persons)[j].Age
			})
		} else if by == OrderByDesc {
			sort.Slice(*persons, func(i, j int) bool {
				return (*persons)[i].Age > (*persons)[j].Age
			})
		}
	} else if field == "Name" {
		if by == OrderByAsc {
			sort.Slice(*persons, func(i, j int) bool {
				return (*persons)[i].Name < (*persons)[j].Name
			})
		} else if by == OrderByDesc {
			sort.Slice(*persons, func(i, j int) bool {
				return (*persons)[i].Name > (*persons)[j].Name
			})
		}
	}
	return nil
}

func WriteSearchErrorResponse(w http.ResponseWriter, response *SearchErrorResponse) {
	w.WriteHeader(http.StatusBadRequest)
	errorText, _ := json.Marshal(response)
	_, _ = w.Write(errorText)
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("AccessToken")
	if accessToken != "MyAccessToken" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var persons []Person
	query := r.FormValue("query")
	orderField := r.FormValue("order_field")
	orderBy := r.FormValue("order_by")
	limitText := r.FormValue("limit")
	offsetText := r.FormValue("offset")
	for _, person := range rootPersons.Persons {
		if strings.Contains(person.Name, query) || strings.Contains(person.About, query) {
			persons = append(persons, person)
		}
	}
	if order, err := strconv.Atoi(orderBy); err == nil {
		if len(orderField) > 0 && orderField != "Name" &&
			orderField != "Id" && orderField != "Age" {
			WriteSearchErrorResponse(w, &SearchErrorResponse{Error: "ErrorBadOrderField"})
			return
		}
		err = SortBy(&persons, orderField, order)
		if err != nil {
			WriteSearchErrorResponse(w, &SearchErrorResponse{Error: "ErrorBadOrderBy"})
			return
		}
	} else {
		WriteSearchErrorResponse(w, &SearchErrorResponse{Error: "ErrorBadOrderBy"})
		return
	}

	limit := len(persons)
	if limitInt, err := strconv.Atoi(limitText); err == nil {
		limit = limitInt
	}
	offset := 0
	if offsetInt, err := strconv.Atoi(offsetText); err == nil {
		offset = offsetInt
	}

	end := offset + limit
	if offset > len(persons) {
		offset = len(persons)
	}
	if end > len(persons) {
		end = len(persons)
	}
	personsText, err := json.Marshal(persons[offset:end])
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	_, _ = w.Write(personsText)
}

func ParseDataset() {
	byteValue, _ := os.ReadFile("dataset.xml")
	err := xml.Unmarshal(byteValue, &rootPersons)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < len(rootPersons.Persons); i++ {
		person := &rootPersons.Persons[i]
		person.Name = fmt.Sprintf("%s %s", person.FirstName, person.LastName)
	}
}

func main() {
	ParseDataset()

	http.HandleFunc("/", SearchServer)
	_ = http.ListenAndServe(":8080", nil)
}
