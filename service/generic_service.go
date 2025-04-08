package service

import (
	"github.com/gin-gonic/gin"
)

type GenericRepository[T any] interface {
	FindAll(query string, filtersQUery map[string]interface{}) ([]T, error)
	FindByID(id string) (T, error)
	Create(item T) error
	Update(id string, item T) error
	Delete(id string) error
}

type GenericService[T any] interface {
	GetAll(query string, filters map[string]interface{}) ([]T, error)
	GetByID(id string) (T, error)
	Create(item T) error
	Update(id string, ctx *gin.Context) error
	Delete(id string) error
}

type GenericServiceImpl[T any] struct {
	Repo GenericRepository[T]
}

func NewGenericService[T any](repo GenericRepository[T]) GenericService[T] {
	return &GenericServiceImpl[T]{Repo: repo}
}

func (s *GenericServiceImpl[T]) GetAll(query string, filters map[string]interface{}) ([]T, error) {
	return s.Repo.FindAll(query, filters)
}

func (s *GenericServiceImpl[T]) GetByID(id string) (T, error) {
	return s.Repo.FindByID(id)
}

func (s *GenericServiceImpl[T]) Create(item T) error {
	return s.Repo.Create(item)
}

func (s *GenericServiceImpl[T]) Update(id string, ctx *gin.Context) error {
	item, err := s.Repo.FindByID(id)

	if err != nil {
		return err
	}

	if err := ctx.ShouldBindJSON(&item); err != nil {
		return err
	}

	return s.Repo.Update(id, item)
}

func (s *GenericServiceImpl[T]) Delete(id string) error {
	return s.Repo.Delete(id)
}
