package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	data := `<root>
<row><id>10</id><first_name>Zoe</first_name><last_name>Zulu</last_name><age>9</age><about>go developer</about><gender>female</gender></row>
<row><id>2</id><first_name>Ann</first_name><last_name>Alpha</last_name><age>30</age><about>designer</about><gender>female</gender></row>
<row><id>30</id><first_name>Bob</first_name><last_name>Beta</last_name><age>20</age><about>go tester</about><gender>male</gender></row>
</root>`
	file := filepath.Join(t.TempDir(), "dataset.xml")
	if err := os.WriteFile(file, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	oldPath := path
	path = file
	defer func() { path = oldPath }()

	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()
	srv := SearchClient{
		URL:         server.URL,
		AccessToken: secretToken,
	}
	zoe := User{
		ID:     10,
		Name:   "Zoe Zulu",
		Age:    9,
		About:  "go developer",
		Gender: "female",
	}
	ann := User{
		ID:     2,
		Name:   "Ann Alpha",
		Age:    30,
		About:  "designer",
		Gender: "female",
	}
	bob := User{
		ID:     30,
		Name:   "Bob Beta",
		Age:    20,
		About:  "go tester",
		Gender: "male",
	}

	tests := []struct {
		name string
		req  SearchRequest
		want []User
		next bool
	}{
		{
			name: "all users",
			req: SearchRequest{
				Limit: 3,
			},
			want: []User{
				zoe,
				ann,
				bob,
			},
		},
		{
			name: "age",
			req: SearchRequest{
				Limit: 3,
				Query: "9",
			},
			want: []User{
				zoe,
			},
		},
		{
			name: "gender",
			req: SearchRequest{
				Limit: 3,
				Query: "female",
			},
			want: []User{
				zoe,
				ann,
			},
		},
		{
			name: "go",
			req: SearchRequest{
				Limit: 3,
				Query: "go",
			},
			want: []User{
				zoe,
				bob,
			},
		},
		{
			name: "name",
			req: SearchRequest{
				Limit: 3,
				Query: "Ann Alpha",
			},
			want: []User{
				ann,
			},
		},
		{
			name: "about",
			req: SearchRequest{
				Limit: 3,
				Query: "go",
			},
			want: []User{
				zoe,
				bob,
			},
		},
		{
			name: "id ascending",
			req: SearchRequest{
				Limit:      3,
				OrderField: "Id",
				OrderBy:    1,
			},
			want: []User{
				ann,
				zoe,
				bob,
			},
		},
		{
			name: "id descending",
			req: SearchRequest{
				Limit:      3,
				OrderField: "Id",
				OrderBy:    -1,
			},
			want: []User{
				bob,
				zoe,
				ann,
			},
		},
		{
			name: "age",
			req: SearchRequest{
				Limit:      3,
				OrderField: "Age",
				OrderBy:    1,
			},
			want: []User{
				zoe,
				bob,
				ann,
			},
		},
		{
			name: "default name sort",
			req: SearchRequest{
				Limit:   3,
				OrderBy: 1,
			},
			want: []User{
				ann,
				bob,
				zoe,
			},
		},
		{
			name: "next page",
			req: SearchRequest{
				Limit: 1,
			},
			want: []User{
				zoe,
			},
			next: true,
		},
		{
			name: "offset after filtering",
			req: SearchRequest{
				Limit:  1,
				Offset: 1,
				Query:  "go",
			},
			want: []User{
				bob,
			},
		},
		{
			name: "offset after sorting",
			req: SearchRequest{
				Limit:   1,
				Offset:  2,
				OrderBy: 1,
			},
			want: []User{
				zoe,
			},
		},
		{
			name: "no matches",
			req: SearchRequest{
				Limit: 3,
				Query: "missing",
			},
			want: []User{},
		},
		{
			name: "offset past end",
			req: SearchRequest{
				Limit:  3,
				Offset: 10,
			},
			want: []User{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := srv.FindUsers(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got.Users, tt.want) || got.NextPage != tt.next {
				t.Errorf("got %+v, want users %+v, next %v", got, tt.want, tt.next)
			}
		})
	}
}

func TestSearchServerParams(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{
			query: "limit=abc",
			want:  "bad param limit:",
		},
		{
			query: "limit=0",
			want:  "bad limit",
		},
		{
			query: "offset=-1",
			want:  "bad offset",
		},
		{
			query: "order_by=2",
			want:  "bad order_by",
		},
		{
			query: "order_field=unknown",
			want:  ErrorBadOrderField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/?"+tt.query, nil)
			req.Header.Set("AccessToken", secretToken)
			rec := httptest.NewRecorder()

			SearchServer(rec, req)

			if rec.Code != 400 {
				t.Fatalf("status = %d, want 400", rec.Code)
			}
			var got SearchErrorResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(got.Error, tt.want) {
				t.Errorf("error = %q, want %q", got.Error, tt.want)
			}
		})
	}
}

func TestSearchServerAccess(t *testing.T) {
	tests := []struct {
		method string
		token  string
		status int
	}{
		{
			method: "GET",
			token:  "wrong",
			status: 401,
		},
		{
			method: "POST",
			token:  secretToken,
			status: 405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			req.Header.Set("AccessToken", tt.token)
			rec := httptest.NewRecorder()

			SearchServer(rec, req)

			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}
		})
	}
}

func TestSearchServerFiles(t *testing.T) {
	badFile := filepath.Join(t.TempDir(), "bad.xml")
	if err := os.WriteFile(badFile, []byte("<root>"), 0600); err != nil {
		t.Fatal(err)
	}
	oldPath := path
	defer func() { path = oldPath }()

	for _, file := range []string{
		filepath.Join(t.TempDir(), "missing.xml"),
		badFile,
	} {
		t.Run(filepath.Base(file), func(t *testing.T) {
			path = file
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("AccessToken", secretToken)
			rec := httptest.NewRecorder()

			SearchServer(rec, req)

			if rec.Code != 500 {
				t.Errorf("status = %d, want 500", rec.Code)
			}
		})
	}
}

func TestFindUsersErrors(t *testing.T) {
	tests := []struct {
		name   string
		req    SearchRequest
		status int
		body   string
		want   string
	}{
		{
			name: "negative limit",
			req: SearchRequest{
				Limit: -1,
			},
			status: 200,
			want:   "limit must be > 0",
		},
		{
			name: "negative offset",
			req: SearchRequest{
				Offset: -1,
			},
			status: 200,
			want:   "offset must be > 0",
		},
		{
			name:   "token",
			status: 401,
			want:   "bad AccessToken",
		},
		{
			name:   "server",
			status: 500,
			want:   "SearchServer fatal error",
		},
		{
			name: "order field",
			req: SearchRequest{
				OrderField: "bad",
			},
			status: 400,
			body:   `{"Error":"OrderField invalid"}`,
			want:   "OrderFeld bad invalid",
		},
		{
			name:   "other error",
			status: 400,
			body:   `{"Error":"bad params"}`,
			want:   "unknown bad request error: bad params",
		},
		{
			name:   "broken error JSON",
			status: 400,
			body:   `{`,
			want:   "cant unpack error json:",
		},
		{
			name:   "broken users JSON",
			status: 200,
			body:   `{`,
			want:   "cant unpack result json:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				if _, err := io.WriteString(w, tt.body); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			srv := SearchClient{
				URL: server.URL,
			}

			got, err := srv.FindUsers(tt.req)
			if got != nil || err == nil || !strings.HasPrefix(err.Error(), tt.want) {
				t.Errorf("got %+v, error %v; want error %q", got, err, tt.want)
			}
		})
	}
}

func TestFindUsersLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer server.Close()
	srv := SearchClient{
		URL:         server.URL,
		AccessToken: secretToken,
	}

	got, err := srv.FindUsers(SearchRequest{
		Limit:      100,
		OrderField: "Id",
		OrderBy:    1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Users) != 25 || !got.NextPage {
		t.Fatalf("got %d users, next %v; want 25, true", len(got.Users), got.NextPage)
	}
	for i, user := range got.Users {
		if user.ID != i {
			t.Errorf("user[%d].ID = %d, want %d", i, user.ID, i)
		}
	}
}
