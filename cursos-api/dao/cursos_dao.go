package dao
type Curso struct {
    CursoID     string  `bson:"curso_id"`
    Nombre      string `bson:"nombre"`
    Descripcion string `bson:"descripcion"`
    CategoriaID int64  `bson:"categoria_id"`
    Capacidad   int64  `bson:"capacidad"`
}

// type Cursos []Curso no se si va aca

