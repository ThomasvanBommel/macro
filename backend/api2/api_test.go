package api2

import (
	"log/slog"
	"macro/db2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
)

func newEngine() (*gin.Engine, *DB) {
	gin.SetMode(gin.TestMode)
	goose.SetLogger(goose.NopLogger())
	return Init(db2.Init(":memory:?cache=shared"), slog.LevelError)
}

func newContext(body string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/api/test", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	return w, c
}

func newRecorder(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	return w
}

func TestInit(t *testing.T) {
	r, db := newEngine()
	defer db.Close()

	assert.NotNil(t, r)
	assert.NotNil(t, db)

	w := newRecorder(r, "GET", "/api/health", "")
	assert.Equal(t, 200, w.Code)
}

func TestBindJSON(t *testing.T) {
	type Input struct {
		Name string `json:"name" binding:"required"`
	}

	_, c := newContext(`{}`)
	err := BindJSON(c, &Input{})
	assert.Error(t, err, "Expected error for missing required field")

	var input Input
	_, c = newContext(`{"name": "test"}`)
	err = BindJSON(c, &input)
	assert.NoError(t, err, "Expected no error for valid JSON")
	assert.Equal(t, "test", input.Name, "Expected Name field to be 'test'")
}
