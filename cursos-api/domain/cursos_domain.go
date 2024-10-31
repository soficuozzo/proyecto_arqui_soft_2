package domain

type CursoData struct {
    CursoID     int64  `json:"curso_id"`
    Nombre      string `json:"nombre"`
    Descripcion string `json:"descripcion"`
    CategoriaID int64  `json:"categoria_id"`
    Capacidad   int64  `json:"capacidad"`
    Imagen      string `json:"imagen,omitempty"`
}


