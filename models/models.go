package models

type TrabajadorInput struct {
	Nombre   string `json:"nombre"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

type ClienteInput struct {
	Nombre          string `json:"nombre"`
	Email           string `json:"email"`
	Telefono        string `json:"telefono"`
	Login           string `json:"login"`
	Password        string `json:"password"`
	CodigoSeguridad string `json:"codigo_seguridad"`
}
