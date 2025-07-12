package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/whosthefunkyy/go-rest-api-example/handlers"
	"github.com/google/go-cmp/cmp"
)

func TestGetUsers(t *testing.T) {
   mockDB, mock, err := sqlmock.New()
   if err != nil {
	t.Fatalf("failed to open mock sql db: %s",err)
   }
   defer mockDB.Close()

   rows:= sqlmock.NewRows([]string{"id","name","age"}).AddRow(1,"Artem",24)
   mock.ExpectQuery("SELECT id, name, age FROM users").WillReturnRows(rows)

   h := &handlers.Handler{DB: mockDB}


   r:= mux.NewRouter()

   r.HandleFunc("/users", h.GetUsers).Methods("GET")
    req, err := http.NewRequest("GET", "/users", nil)
    if err != nil {
        t.Fatalf("failed to create request: %s", err)
    }

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rr.Code)
	}

	expected:= `[{"id":1,"name":"Artem","age":24}]`

	 got := strings.TrimSpace(rr.Body.String())

	expected = strings.TrimSpace(expected)


	if expected != got {
		t.Errorf("got: but was needed: %s, %s",got, expected)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}
func TestGetUserByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql db: %s", err)

	}
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow(1, "Artem", 24)
	mock.ExpectQuery("SELECT id,name,age FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	h := &handlers.Handler{DB: mockDB}
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", h.GetUser).Methods("GET")

	req,err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr,req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", rr.Code)
	}

	expected := `{
		"id": 1,
		"name": "Artem",
		"age": 24,
		"_links": {
			"self": "/api/v1/users/1",
			"update": "/api/v1/users/1",
			"delete": "/api/v1/users/1",
			"allUsers": "/api/v1/users"
		}
	}`
	var expectedMap map[string]interface{}
	var gotMap map[string]interface{}

	if err := json.Unmarshal([]byte(expected), &expectedMap); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %s", err)
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &gotMap); err != nil {
		t.Errorf("failed to unmarshal response json %s", err)
	}

	if diff := cmp.Diff(expectedMap, gotMap); diff != "" {
	t.Errorf("response body mismatch (-want +got):\n%s", diff)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %s", err)
	}
}
 func TestCreateUser(t *testing.T){
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock sql db: %s", err)
	}
	defer mockDB.Close()

	mock.ExpectQuery("INSERT INTO users \\(name, age\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
					WithArgs("Artem",24).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	h:= &handlers.Handler{DB: mockDB}
	r:= mux.NewRouter()
	r.HandleFunc("/users", h.CreateUser).Methods("POST")

	body := `{"name": "Artem", "age": 24}`

	req, err := http.NewRequest("POST", "/users", strings.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	rr:= httptest.NewRecorder()
	r.ServeHTTP(rr,req)

	if rr.Code != http.StatusCreated {
	t.Errorf("expected 201 Created, got %d", rr.Code)
	}
	expected := `{
		"id": 1,
		"name": "Artem",
		"age": 24,
		"_links": {
			"self": "/api/v1/users/1",
			"update": "/api/v1/users/1",
			"delete": "/api/v1/users/1",
			"allUsers": "/api/v1/users"
		}
	}`
	var expectedMap map[string]interface{}
	var gotMap map[string]interface{}

	if err := json.Unmarshal([]byte(expected), &expectedMap); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %s", err)
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &gotMap); err != nil {
		t.Fatalf("failed to unmarshal response JSON: %s", err)
	}
	if diff := cmp.Diff(expectedMap, gotMap); diff != "" {
	t.Errorf("response body mismatch (-want +got):\n%s", diff)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %s", err)
	}

 }