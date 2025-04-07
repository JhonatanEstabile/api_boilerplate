package service

import "github.com/gin-gonic/gin"

type GenericRepository[T any] interface {
	FindAll() ([]T, error)
	FindByID(id int64) (T, error)
	Create(item T) error
	Update(id int64, item T) error
	Delete(id int64) error
}

type GenericService[T any] interface {
	GetAll() ([]T, error)
	GetByID(id int64) (T, error)
	Create(item T) error
	Update(id int64, ctx *gin.Context) error
	Delete(id int64) error
}

type GenericServiceImpl[T any] struct {
	Repo GenericRepository[T]
}

func NewGenericService[T any](repo GenericRepository[T]) GenericService[T] {
	return &GenericServiceImpl[T]{Repo: repo}
}

func (s *GenericServiceImpl[T]) GetAll() ([]T, error) {
	return s.Repo.FindAll()
}

func (s *GenericServiceImpl[T]) GetByID(id int64) (T, error) {
	return s.Repo.FindByID(id)
}

func (s *GenericServiceImpl[T]) Create(item T) error {
	return s.Repo.Create(item)
}

func (s *GenericServiceImpl[T]) Update(id int64, ctx *gin.Context) error {
	item, err := s.Repo.FindByID(id)

	if err != nil {
		return err
	}

	if err := ctx.ShouldBindJSON(&item); err != nil {
		return err
	}

	return s.Repo.Update(id, item)
}

func (s *GenericServiceImpl[T]) Delete(id int64) error {
	return s.Repo.Delete(id)
}
