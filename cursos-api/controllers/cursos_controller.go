package controllers

import (
    "net/http"
    "strconv"
    "proyecto_arqui_soft_2/cursos-api/domain"
    "proyecto_arqui_soft_2/cursos-api/services"
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"
)


func CrearCurso(c *gin.Context) {
    var nuevoCurso domain.CursoData
    if err := c.ShouldBindJSON(&nuevoCurso); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := services.CrearCurso(nuevoCurso); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Curso creado exitosamente"})
}
func GetCursoById(c *gin.Context) {
    cursoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    curso, err := services.GetCursoByID(cursoID)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, curso)
}


func GetCursos(c *gin.Context) {
    cursos, err := services.GetAllCursos()
    if err != nil || len(cursos) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No hay cursos disponibles"})
        return
    }

    c.JSON(http.StatusOK, cursos)
}

func DeleteCourseName(c *gin.Context) {
    var json struct {
        Nombre string `json:"nombre" binding:"required"`
    }

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := services.DeleteCursoByName(json.Nombre); err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Curso no encontrado"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Curso eliminado"})
}

func EditCourse(c *gin.Context) {
    cursoID, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    var updateData map[string]interface{}
    if err := c.ShouldBindJSON(&updateData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := services.UpdateCurso(cursoID, updateData); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Curso actualizado exitosamente"})
}
