package main

import (
	"sort"
	"strings"
)

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
