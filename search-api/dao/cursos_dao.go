package dao
type Curso struct {
	CursoID     string `bson:"_id,omitempty"`
	Nombre      string `bson:"nombre"`
	Descripcion string `bson:"descripcion"`
	Categoria   string `bson:"categoria"`
	Capacidad   int64  `bson:"capacidad"`
	Requisito   string  `bson:"requisito"`
	Duracion   int64  `bson:"duracion"`
	Imagen   string  `bson:"imagen"`
	Valoracion   int64  `bson:"valoracion"`
	Profesor   string  `bson:"profesor"`
}

type CursosData []Curso