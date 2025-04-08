package controller

import (
	"net/http"

	"api_boilerplate/middleware"
	"api_boilerplate/service"

	"github.com/gin-gonic/gin"
)

type GenericController[T any] struct {
	Service service.GenericService[T]
}

func NewGenericController[T any](s service.GenericService[T]) *GenericController[T] {
	return &GenericController[T]{Service: s}
}

func (c *GenericController[T]) RegisterRoutes(r *gin.Engine, path string) {
	group := r.Group(path)
	group.GET("/", middleware.FilterMiddleware(), c.GetAll)
	group.GET("/:id", c.GetByID)
	group.POST("/", c.Create)
	group.PUT("/:id", c.Update)
	group.DELETE("/:id", c.Delete)
}

func (c *GenericController[T]) GetAll(ctx *gin.Context) {
	query, _ := ctx.Get("filtersSQL")
	params, _ := ctx.Get("filtersQuery")

	items, err := c.Service.GetAll(query.(string), params.(map[string]interface{}))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (c *GenericController[T]) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	item, err := c.Service.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *GenericController[T]) Create(ctx *gin.Context) {
	var item T

	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.Service.Create(item); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *GenericController[T]) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.Service.Update(id, ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *GenericController[T]) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.Service.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
