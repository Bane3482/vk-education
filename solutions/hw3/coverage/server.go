package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

const secretToken = "access_token"

var path = "dataset.xml"

var (
	ErrBadLimit      = errors.New("bad limit")
	ErrBadOffset     = errors.New("bad offset")
	ErrBadOrderField = errors.New("bad order_field")
	ErrBadOrderBy    = errors.New("bad order_by")
)

type Root struct {
	Version string `xml:"version,attr"`
	Rows    []Row  `xml:"row"`
}

type Row struct {
	ID        int    `xml:"id"`
	Age       int    `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
}

type SearchParams struct {
	Limit      int
	Offset     int
	Query      string
	OrderField string
	OrderBy    int
}

func (s SearchParams) Validate() error {
	if s.Limit <= 0 {
		return ErrBadLimit
	}
	if s.Offset < 0 {
		return ErrBadOffset
	}
	if s.OrderField != "Id" && s.OrderField != "Age" && s.OrderField != "Name" {
		return ErrBadOrderField
	}
	if s.OrderBy != OrderByAsIs && s.OrderBy != OrderByAsc && s.OrderBy != OrderByDesc {
		return ErrBadOrderBy
	}
	return nil
}

func parseParams(r *http.Request) (SearchParams, error) {
	values := r.URL.Query()
	params := SearchParams{Limit: 10, Query: values.Get("query"), OrderField: "Name"}

	numParamsList := map[string]*int{
		"limit":    &params.Limit,
		"offset":   &params.Offset,
		"order_by": &params.OrderBy,
	}

	for name, target := range numParamsList {
		if value := values.Get(name); value != "" {
			number, err := strconv.Atoi(value)
			if err != nil {
				return SearchParams{}, fmt.Errorf("bad param %s: %w", name, err)
			}
			*target = number
		}
	}

	switch field := values.Get("order_field"); strings.ToLower(field) {
	case "", "name":
		params.OrderField = "Name"
	case "id":
		params.OrderField = "Id"
	case "age":
		params.OrderField = "Age"
	default:
		params.OrderField = field
	}
	return params, nil
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write response: %v", err)
	}
}

func SearchUsers(rows []Row, params SearchParams) []User {
	users := make([]User, 0)
	for _, row := range rows {
		name := row.FirstName + " " + row.LastName
		if strings.Contains(name, params.Query) || strings.Contains(row.About, params.Query) {
			users = append(users, User{
				ID:     row.ID,
				Name:   name,
				Age:    row.Age,
				Gender: row.Gender,
				About:  row.About,
			})
		}
	}

	if params.OrderBy != OrderByAsIs {
		sort.Slice(users, func(i, j int) bool {
			if params.OrderBy == OrderByDesc {
				i, j = j, i
			}
			switch params.OrderField {
			case "Id":
				return users[i].ID < users[j].ID
			case "Age":
				return users[i].Age < users[j].Age
			default:
				return users[i].Name < users[j].Name
			}
		})
	}

	if params.Offset >= len(users) {
		return []User{}
	}
	users = users[params.Offset:]
	if len(users) > params.Limit {
		users = users[:params.Limit]
	}
	return users
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, SearchErrorResponse{Error: "method not allowed"})
		return
	}
	if r.Header.Get("AccessToken") != secretToken {
		writeJSON(w, http.StatusUnauthorized, SearchErrorResponse{Error: "bad access token"})
		return
	}

	params, err := parseParams(r)
	if err == nil {
		err = params.Validate()
	}
	if err != nil {
		message := err.Error()
		if errors.Is(err, ErrBadOrderField) {
			message = ErrorBadOrderField
		}
		writeJSON(w, http.StatusBadRequest, SearchErrorResponse{Error: message})
		return
	}

	file, err := os.Open(path)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, SearchErrorResponse{Error: "cannot open dataset"})
		return
	}
	defer file.Close()

	var root Root
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		writeJSON(w, http.StatusInternalServerError, SearchErrorResponse{Error: "cannot decode dataset"})
		return
	}

	writeJSON(w, http.StatusOK, SearchUsers(root.Rows, params))
}
