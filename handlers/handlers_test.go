package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"context"
	"time"
)

// ---- GET /users ---- 1
func TestGetUsers_Success(t *testing.T) {
	h := &Handler{Repo: &mockUserRepository{}}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetAllUsers(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var users []map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Artem", users[0]["name"])
}

// ---- GET /users ---- 2

func TestGetUsers_DBError(t *testing.T) {
	handler := &Handler{Repo: &mockFailingUserRepository{}}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.GetAllUsers(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ---- GET BY ID /users/id ---- 1
func TestGetUser_Success(t *testing.T) {
	handler := &Handler{Repo: &mockUserRepository{}}

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetUserByID(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var user map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, "Artem", user["name"])
}

// ---- GET BY ID /users/id ---- 2
func TestGetUser_NotFound(t *testing.T) {
	handler := &Handler{Repo: &mockUserRepository{}}

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetUserByID(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// ---- GET BY ID /users/id ---- 3
func TestGetUser_InvalidID(t *testing.T) {
	handler := &Handler{Repo: &mockUserRepository{}}

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.GetUserByID(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// ---- GET BY ID /users/id ---- 4
func TestGetUser_DBError(t *testing.T) {
	handler := &Handler{Repo: &mockErrorUserRepository{}}

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetUserByID(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ---- GET BY ID /users/id ---- 5
func TestGetUser_Timeout(t *testing.T) {
	handler := &Handler{Repo: &mockUserRepository{}}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	time.Sleep(2 * time.Nanosecond)
	cancel()

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil).WithContext(ctx)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetUserByID(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusRequestTimeout, resp.StatusCode)
}

