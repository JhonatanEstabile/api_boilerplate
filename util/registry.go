package util

import (
	"api_boilerplate/controller"
	"api_boilerplate/model"
	"api_boilerplate/repository"
	"api_boilerplate/service"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RegisterGenericResource[T any](r *gin.Engine, db *sqlx.DB, path string, fields []string) {
	repo := repository.NewSqlxRepository[T](db, path, fields)
	service := service.NewGenericService(repo)
	controller := controller.NewGenericController(service)
	controller.RegisterRoutes(r, "/"+path)
}

func RegisterDomains(r *gin.Engine, db *sqlx.DB) {
	RegisterGenericResource[model.User](r, db, "user", model.UserFields)
	RegisterGenericResource[model.Product](r, db, "product", model.ProductFields)
	RegisterGenericResource[model.Store](r, db, "store", model.StoreFields)
}
