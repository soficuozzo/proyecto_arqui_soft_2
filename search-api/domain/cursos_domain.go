package domain

type CursoData struct {
	CursoID     string `json:"curso_id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Categoria   string `json:"categoria"`
	Capacidad   int64  `json:"capacidad"`
}


type CursoNew struct {
	Operation string `json:"operation"`
	CursoID string `json:"curso_id"`
}