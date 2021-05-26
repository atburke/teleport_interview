package main

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "net/http"
  "net/http/httptest"
)

func TestPing(t *testing.T) {
  router := setupRouter()
  w := httptest.NewRecorder()
  request := httptest.NewRequest("GET", "/ping", nil)
  router.ServeHTTP(w, request)

  assert.Equal(t, http.StatusOK, w.Code)
  assert.Equal(t, "pong", w.Body.String())
}
