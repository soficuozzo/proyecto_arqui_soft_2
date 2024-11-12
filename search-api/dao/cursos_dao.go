package dao
type Curso struct {
	CursoID     string `bson:"_id,omitempty"`
	Nombre      string `bson:"nombre"`
	Descripcion string `bson:"descripcion"`
	Categoria   string `bson:"categoria"`
	Capacidad   int64  `bson:"capacidad"`
}

type CursosData []Curso