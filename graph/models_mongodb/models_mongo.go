package models_mongodb

// eSTO representa a la collection obra ultimate
type ObraMongo struct {
	ID            int32        `bson:"_id"`
	Nombre        string       `bson:"nombre"`
	IDArtista     int32        `bson:"id_artista"`
	IDGenero      int32        `bson:"id_genero"`
	PrecioObra    int32        `bson:"precio_obra"`
	FechaCreacion string       `bson:"fecha_creacion"`
	Status        string       `bson:"status"`
	Foto          string       `bson:"foto"`
	Artista       ArtistaMongo `bson:"artista"` 
	Genero        GeneroMongo  `bson:"genero"`  
}

// Esto representa al objeto embebido en obra ultimata
type ArtistaMongo struct {
	ID              int32  `bson:"_id"`
	Nombre          string `bson:"nombre"`
	FechaNacimiento string `bson:"fecha_nacimiento"`
	Nacionalidad    string `bson:"nacionalidad"`
	Biografia       string `bson:"biografia"`
	Foto            string `bson:"foto"`
}

// Lo mismo que el anterior pero para genero
type GeneroMongo struct {
	ID     int32  `bson:"_id"`
	Nombre string `bson:"nombre"`
}