package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockRepository é um mock simples para testar o service
type MockRepository[T any] struct {
	FindAllFn  func() ([]T, error)
	FindByIDFn func(int64) (T, error)
	CreateFn   func(T) error
	UpdateFn   func(int64, T) error
	DeleteFn   func(int64) error
}

func (m *MockRepository[T]) FindAll() ([]T, error) {
	return m.FindAllFn()
}

func (m *MockRepository[T]) FindByID(id int64) (T, error) {
	return m.FindByIDFn(id)
}

func (m *MockRepository[T]) Create(item T) error {
	return m.CreateFn(item)
}

func (m *MockRepository[T]) Update(id int64, item T) error {
	return m.UpdateFn(id, item)
}

func (m *MockRepository[T]) Delete(id int64) error {
	return m.DeleteFn(id)
}

type TestModel struct {
	ID   int64
	Name string
}

func TestGenericService_GetAll(t *testing.T) {
	mockRepo := &MockRepository[TestModel]{
		FindAllFn: func() ([]TestModel, error) {
			return []TestModel{{ID: 1, Name: "Test"}}, nil
		},
	}
	service := NewGenericService[TestModel](mockRepo)
	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ID)
}

func TestGenericService_GetByID(t *testing.T) {
	mockRepo := &MockRepository[TestModel]{
		FindByIDFn: func(id int64) (TestModel, error) {
			return TestModel{ID: id, Name: "Test"}, nil
		},
	}
	service := NewGenericService[TestModel](mockRepo)
	result, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
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
	err := service.Create(TestModel{ID: 2, Name: "New"})

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestGenericService_Update(t *testing.T) {
	called := false
	mockRepo := &MockRepository[TestModel]{
		UpdateFn: func(id int64, item TestModel) error {
			called = true
			return nil
		},
	}
	service := NewGenericService[TestModel](mockRepo)
	err := service.Update(1, TestModel{ID: 1, Name: "Updated"})

	assert.NoError(t, err)
	assert.True(t, called)
}

func TestGenericService_Delete(t *testing.T) {
	called := false
	mockRepo := &MockRepository[TestModel]{
		DeleteFn: func(id int64) error {
			called = true
			return nil
		},
	}
	service := NewGenericService[TestModel](mockRepo)
	err := service.Delete(1)

	assert.NoError(t, err)
	assert.True(t, called)
}
