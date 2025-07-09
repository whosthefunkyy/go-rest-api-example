package handlers_test

import (
	
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/whosthefunkyy/go-rest-api-example/handlers"
)

func TestGetUsers(t *testing.T) {
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open mock sql db %s", err)
    }

    defer mockDB.Close()

    rows := sqlmock.NewRows([]string{"id","name","age"}).AddRow(1,"Artem", 24)
    mock.ExpectQuery("SELECT id, name, age FROM users").WillReturnRows(rows)

    h := &handlers.Handler{DB: mockDB}

    r:= mux.NewRouter()

    r.HandleFunc("/users", h.GetUsers).Methods("GET")

    req,err := http.NewRequest("GET", "/users", nil)
    if err != nil {
        t.Fatalf("failed to create request %s", err)
    }

    rr := httptest.NewRecorder()

    r.ServeHTTP(rr,req)

    if rr.Code != http.StatusOK {
         t.Errorf("expected status 200 OK, got %d", rr.Code)
    }

    expected := `[{"id":1,"name":"Artem","age":24}]`

    got := rr.Body.String()

    got = strings.TrimSpace(got)
    expected = strings.TrimSpace(expected)  

    if got != expected {
         t.Errorf("unexpected body:\n got: %s\nwant: %s", got, expected)
    }


    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectetion: %s", err)
    }
}