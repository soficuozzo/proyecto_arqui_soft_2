package services

import (
    "context"
    "proyecto_arqui_soft_2/cursos-api/dao"
    "proyecto_arqui_soft_2/cursos-api/domain"
    "proyecto_arqui_soft_2/cursos-api/clients"
    "time"
    "go.mongodb.org/mongo-driver/bson"
)

func ToCursoDAO(curso domain.CursoData) dao.Curso {
    return dao.Curso{
        CursoID:     curso.CursoID,
        Nombre:      curso.Nombre,
        Descripcion: curso.Descripcion,
        CategoriaID: curso.CategoriaID,
        Capacidad:   curso.Capacidad,
    }
}

func ToDomainCurso(curso dao.Curso) domain.CursoData {
    return domain.CursoData{
        CursoID:     curso.CursoID,
        Nombre:      curso.Nombre,
        Descripcion: curso.Descripcion,
        CategoriaID: curso.CategoriaID,
        Capacidad:   curso.Capacidad,
    }
}

func CrearCurso(curso domain.CursoData) error {
    collection := GetCursoCollection() 
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursoDAO := ToCursoDAO(curso)  // Convierte la estructura a la versión de Mongo
    _, err := collection.InsertOne(ctx, cursoDAO)
    if err != nil {
        return err
    }

/*     // Notificar a RabbitMQ
    err = clients.PublishCursoCreado(cursoDAO)
    if err != nil {
        return err
    } */

    return nil
}

func GetCursoByID(cursoID int64) (domain.CursoData, error) {
    var cursoDAO dao.Curso
    collection := GetCursoCollection()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := collection.FindOne(ctx, bson.M{"curso_id": cursoID}).Decode(&cursoDAO)
    if err != nil {
        return domain.CursoData{}, err
    }

    return ToDomainCurso(cursoDAO), nil  
}


func GetAllCursos() ([]domain.CursoData, error) {
    var cursosDAO []dao.Curso
    collection := GetCursoCollection()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }

    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var cursoDAO dao.Curso
        err := cursor.Decode(&cursoDAO)
        if err != nil {
            return nil, err
        }
        cursosDAO = append(cursosDAO, cursoDAO)
    }

    var cursos []domain.CursoData
    for _, cursoDAO := range cursosDAO {
        cursos = append(cursos, ToDomainCurso(cursoDAO))
    }

    return cursos, nil
}


func DeleteCursoByName(nombre string) error {
    collection := GetCursoCollection()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := collection.DeleteOne(ctx, bson.M{"nombre": nombre})
    if err != nil {
        return err
    }

    if result.DeletedCount == 0 {
        return mongo.ErrNoDocuments
    }

    return nil
}


func UpdateCurso(cursoID int64, updateData bson.M) error {
    collection := GetCursoCollection()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{"curso_id": cursoID}
    update := bson.M{"$set": updateData}

    _, err := collection.UpdateOne(ctx, filter, update)
    return err
}

