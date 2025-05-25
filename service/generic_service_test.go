package service

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRepository[T any] struct {
	FindAllFn  func() ([]T, error)
	FindByIDFn func(string) (T, error)
	CreateFn   func(T) error
	UpdateFn   func(string, T) error
	DeleteFn   func(string) error
}

func (m *MockRepository[T]) FindAll(query string, filters map[string]interface{}) ([]T, error) {
	return m.FindAllFn()
}
func (m *MockRepository[T]) FindByID(id string) (T, error)  { return m.FindByIDFn(id) }
func (m *MockRepository[T]) Create(item T) error            { return m.CreateFn(item) }
func (m *MockRepository[T]) Update(id string, item T) error { return m.UpdateFn(id, item) }
func (m *MockRepository[T]) Delete(id string) error         { return m.DeleteFn(id) }

type TestModel struct {
	ID   string
	Name string
}

func TestGenericService_GetAll(t *testing.T) {
	mockRepo := &MockRepository[TestModel]{
		FindAllFn: func() ([]TestModel, error) {
			return []TestModel{{ID: "01JW4MH8S671QVVGD0NYY1XWAP", Name: "Test"}}, nil
		},
	}

	filters := make(map[string]interface{})
	filters["id"] = "01JW4MH8S671QVVGD0NYY1XWAP"
	query := "WHERE id = :id"

	service := NewGenericService[TestModel](mockRepo)
	result, err := service.GetAll(query, filters)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "01JW4MH8S671QVVGD0NYY1XWAP", result[0].ID)
}

func TestGenericService_GetByID(t *testing.T) {
	mockRepo := &MockRepository[TestModel]{
		FindByIDFn: func(id string) (TestModel, error) {
			return TestModel{ID: id, Name: "Test"}, nil
		},
	}
	service := NewGenericService[TestModel](mockRepo)
	result, err := service.GetByID("01JW4MH8S671QVVGD0NYY1XWAP")

	assert.NoError(t, err)
	assert.Equal(t, "01JW4MH8S671QVVGD0NYY1XWAP", result.ID)
}

func TestGenericService_Create(t *testing.T) {
	called := false
	mockRepo := &MockRepository[TestModel]{
		CreateFn: func(item TestModel) error {
			called = true
			return nil
		},
	}
	service := NewGenericService[TestModel](mockRepo)
	err := service.Create(TestModel{ID: "01JW4MH8S671QVVGD0NYY1XWAP", Name: "New"})

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestGenericService_Delete(t *testing.T) {
	called := false
	mockRepo := &MockRepository[TestModel]{
		DeleteFn: func(id string) error {
			called = true
			return nil
		},
	}
	service := NewGenericService[TestModel](mockRepo)
	err := service.Delete("01JW4MH8S671QVVGD0NYY1XWAP")

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestGenericService_Update(t *testing.T) {
	mockModel := TestModel{ID: "01JW4MH8S671QVVGD0NYY1XWAP", Name: "Old Name"}
	mockRepo := &MockRepository[TestModel]{
		FindByIDFn: func(id string) (TestModel, error) {
			return mockModel, nil
		},
		UpdateFn: func(id string, updated TestModel) error {
			assert.Equal(t, "01JW4MH8S671QVVGD0NYY1XWAP", id)
			assert.Equal(t, "New Name", updated.Name)
			return nil
		},
	}
	service := NewGenericService[TestModel](mockRepo)

	// Criando contexto simulado com JSON
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	body := bytes.NewBufferString(`{"name":"New Name"}`)
	ctx.Request, _ = http.NewRequest(http.MethodPut, "/test/1", body)
	ctx.Request.Header.Set("Content-Type", "application/json")

	err := service.Update("01JW4MH8S671QVVGD0NYY1XWAP", ctx)

	assert.NoError(t, err)
}
