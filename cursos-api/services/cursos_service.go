package services

import (
	"context"
	"fmt"
	"proyecto_arqui_soft_2/cursos-api/dao"
	"proyecto_arqui_soft_2/cursos-api/domain"

	"go.mongodb.org/mongo-driver/bson" // Asegúrate de importar el paquete bson
)

type Repository interface {
	GetCursoByID(ctx context.Context, id string) (dao.Curso, error)
	Create(ctx context.Context, curso dao.Curso) (string, error)
	Update(ctx context.Context, id string, updateData bson.M) error // Cambiado a bson.M
	Delete(ctx context.Context, id string) error
}

type CursoService struct {
	mainRepository Repository
}

func NewCursoService(mainRepository Repository) CursoService {
	return CursoService{
		mainRepository: mainRepository,
	}
}

func (service CursoService) GetCursoByID(ctx context.Context, id string) (domain.CursoData, error) {
	cursoDAO, err := service.mainRepository.GetCursoByID(ctx, id)
	if err != nil {
		return domain.CursoData{}, fmt.Errorf("error obteniendo curso: %w", err)
	}

	return domain.CursoData{
		CursoID:     cursoDAO.CursoID,
		Nombre:      cursoDAO.Nombre,
		Descripcion: cursoDAO.Descripcion,
		Categoria:   cursoDAO.Categoria,
		Capacidad:   cursoDAO.Capacidad,
	}, nil
}

func (service CursoService) Create(ctx context.Context, curso domain.CursoData) (string, error) {
	cursoDAO := dao.Curso{
		Nombre:      curso.Nombre,
		Descripcion: curso.Descripcion,
		Categoria:   curso.Categoria,
		Capacidad:   curso.Capacidad,
	}

	id, err := service.mainRepository.Create(ctx, cursoDAO)
	if err != nil {
		return "", fmt.Errorf("error creando curso: %w", err)
	}

	return id, nil
}

func (service CursoService) Update(ctx context.Context, curso domain.CursoData) error {
	updateData := bson.M{}
	if curso.Nombre != "" {
		updateData["nombre"] = curso.Nombre
	}
	if curso.Descripcion != "" {
		updateData["descripcion"] = curso.Descripcion
	}
	if curso.Categoria != "" {
		updateData["categoria_id"] = curso.Categoria
	}
	if curso.Capacidad != 0 {
		updateData["capacidad"] = curso.Capacidad
	}

	if len(updateData) == 0 {
		return fmt.Errorf("no hay campos para actualizar")
	}

	err := service.mainRepository.Update(ctx, curso.CursoID, updateData) // Se queda igual
	if err != nil {
		return fmt.Errorf("error actualizando curso: %w", err)
	}

	return nil
}

func (service CursoService) Delete(ctx context.Context, id string) error {
	err := service.mainRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error eliminando curso: %w", err)
	}
	return nil
}
