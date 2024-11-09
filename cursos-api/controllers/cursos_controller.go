package controllers

import (
	"context"
	"fmt"
	"net/http"
	"proyecto_arqui_soft_2/cursos-api/domain"
	"strings"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetCursoByID(ctx context.Context, id string) (domain.CursoData, error)
	Create(ctx context.Context, curso domain.CursoData) (string, error)
	Update(ctx context.Context, curso domain.CursoData) error
	Delete(ctx context.Context, id string) error
}

type CursoController struct {
	service Service
}

func NewCursoController(service Service) CursoController {
	return CursoController{
		service: service,
	}
}

func (controller CursoController) GetCursoByID(ctx *gin.Context) {
	cursoID := strings.TrimSpace(ctx.Param("id"))

	curso, err := controller.service.GetCursoByID(ctx.Request.Context(), cursoID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("error obteniendo curso: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, curso)
}

func (controller CursoController) Create(ctx *gin.Context) {

	var curso domain.CursoData
	if err := ctx.ShouldBindJSON(&curso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("solicitud inválida: %s", err.Error()),
		})
		return
	}

	id, err := controller.service.Create(ctx.Request.Context(), curso)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error creando curso: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (controller CursoController) Update(ctx *gin.Context) {

	id := strings.TrimSpace(ctx.Param("id"))

	var curso domain.CursoData
	if err := ctx.ShouldBindJSON(&curso); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("solicitud inválida: %s", err.Error()),
		})
		return
	}

	curso.CursoID = id

	if err := controller.service.Update(ctx.Request.Context(), curso); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error actualizando curso: %s", err.Error()),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Curso actualizado",
	})
}

func (controller CursoController) Delete(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if err := controller.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error eliminando curso: %s", err.Error()),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Curso eliminado",
	})
}
