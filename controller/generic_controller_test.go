package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type MockService[T any] struct {
	GetAllFn  func() ([]T, error)
	GetByIDFn func(int64) (T, error)
	CreateFn  func(T) error
	UpdateFn  func(int64, *gin.Context) error
	DeleteFn  func(int64) error
}

func (m *MockService[T]) GetAll(query string, filters map[string]interface{}) ([]T, error) {
	return m.GetAllFn()
}
func (m *MockService[T]) GetByID(id int64) (T, error)             { return m.GetByIDFn(id) }
func (m *MockService[T]) Create(item T) error                     { return m.CreateFn(item) }
func (m *MockService[T]) Update(id int64, ctx *gin.Context) error { return m.UpdateFn(id, ctx) }
func (m *MockService[T]) Delete(id int64) error                   { return m.DeleteFn(id) }

func setupRouter[T any](controller *GenericController[T]) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	controller.RegisterRoutes(r, "/test")
	return r
}

func TestGenericController_GetAll(t *testing.T) {
	service := &MockService[TestModel]{
		GetAllFn: func() ([]TestModel, error) {
			return []TestModel{{ID: 1, Name: "Test"}}, nil
		},
	}
	ctrl := NewGenericController(service)
	router := setupRouter(ctrl)

	req, _ := http.NewRequest("GET", "/test/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	var body []TestModel
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), body[0].ID)
}

func TestGenericController_Create(t *testing.T) {
	called := false
	service := &MockService[TestModel]{
		CreateFn: func(item TestModel) error {
			called = true
			return nil
		},
	}
	ctrl := NewGenericController(service)
	router := setupRouter(ctrl)

	body, _ := json.Marshal(TestModel{Name: "New"})
	req, _ := http.NewRequest("POST", "/test/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 201, resp.Code)
	assert.True(t, called)
}

func TestGenericController_Delete(t *testing.T) {
	called := false
	service := &MockService[TestModel]{
		DeleteFn: func(id int64) error {
			called = true
			return nil
		},
	}
	ctrl := NewGenericController(service)
	router := setupRouter(ctrl)

	req, _ := http.NewRequest("DELETE", "/test/1", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 204, resp.Code)
	assert.True(t, called)
}

func TestGenericController_Update(t *testing.T) {
	called := false
	service := &MockService[TestModel]{
		UpdateFn: func(id int64, ctx *gin.Context) error {
			called = true
			return nil
		},
	}
	ctrl := NewGenericController(service)
	router := setupRouter(ctrl)

	body, _ := json.Marshal(TestModel{Name: "Updated"})
	req, _ := http.NewRequest("PUT", "/test/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, 200, resp.Code)
	assert.True(t, called)
}
