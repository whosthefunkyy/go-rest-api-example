package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	
	"github.com/stretchr/testify/assert"
)

// ---- GET /users ----

func TestGetUsers_Success(t *testing.T) {
	h := &Handler{Repo: &mockUserRepository{}}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	h.GetUsers(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var users []map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Artem", users[0]["name"])
}
func TestGetUsers_DBError(t *testing.T) {
	handler := &Handler{Repo: &mockFailingUserRepository{}}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.GetUsers(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}