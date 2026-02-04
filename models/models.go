package models

type Trabajador struct {
	Id       int64  `json:"id_trabajador"`
	Nombre   string `json:"nombre"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

type Cliente struct {
	Id              int64  `json:"id"`
	Nombre          string `json:"nombre"`
	Email           string `json:"email"`
	Telefono        string `json:"telefono"`
	Login           string `json:"login"`
	Password        string `json:"password"`
	CodigoSeguridad string `json:"codigo_seguridad"`
}

type Artista struct {
	Id_artista       int64  `json:"id_artista"`
	Nombre           string `json:"nombre"`
	Fecha_nacimiento string `json:"fecha_nacimiento"`
	Nacionalidad     string `json:"nacionalidad"`
	Biografia        string `json:"biografia"`
	Foto             string `json:"foto"`
}

type Genero struct {
	Id_genero int64  `json:"id_genero"`
	Nombre    string `json:"nombre"`
}

type Artista_Genero struct {
	Id_artista int64 `json:"id_artista"`
	Id_genero  int64 `json:"id_genero"`
}

type Obra struct {
	Id_obra        int64   `json:"id_obra"`
	Nombre         string  `json:"nombre"`
	Id_artista     int64   `json:"id_artista"`
	Id_genero      int64   `json:"id_genero"`
	Precio         float64 `json:"precio"`
	Fecha_creacion string  `json:"fecha_creacion"`
	Estatus        string  `json:"estatus"`
	Foto           string  `json:"foto"`
	Material       string  `json:"material"`
	Peso           float64 `json:"peso"`
	Dimensiones    string  `json:"dimensiones"`
}

type Orden struct {
	Id_orden       int64   `json:"id_orden"`
	Id_cliente     int64   `json:"id_cliente"`
	Id_obra        int64   `json:"id_obra"`
	Id_trabajador  int64   `json:"id_trabajador"`
	Precio         float64 `json:"precio"`
	Iva            float64 `json:"iva"`
	Ganancia_museo float64 `json:"ganancia_museo"`
	Total          float64 `json:"total"`
	Fecha_orden    string  `json:"fecha_orden"`
	Estatus        string  `json:"estatus"`
}
