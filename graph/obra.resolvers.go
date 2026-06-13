package graph

import (
	"context"
	"fmt"
	"io"
	"log"
	"main/data_bases/mongodb"
	"main/graph/model"
	"main/graph/models_mongodb"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Como en mongo no existe ENUM hay que ingeniarselas
var validStatuses = map[string]bool{"DISPONIBLE": true, "RESERVADA": true, "VENDIDA": true}

// Helper clave para no repetir codigo, mapea la estructura de base de datos directamente al modelo de GraphQL
func toGraphQLObra(mo *models_mongodb.ObraMongo) *model.Obra {
	return &model.Obra{
		ID:            fmt.Sprintf("%d", mo.ID),
		Nombre:        mo.Nombre,
		IDArtista:     fmt.Sprintf("%d", mo.IDArtista),
		IDGenero:      fmt.Sprintf("%d", mo.IDGenero),
		Precio:        mo.PrecioObra,
		FechaCreacion: mo.FechaCreacion,
		Status:        model.StatusObra(mo.Status),
		Foto:          mo.Foto,
		Artista: &model.Artista{
			ID:              fmt.Sprintf("%d", mo.Artista.ID),
			Nombre:          mo.Artista.Nombre,
			FechaNacimiento: mo.Artista.FechaNacimiento,
			Nacionalidad:    mo.Artista.Nacionalidad,
			Biografia:       mo.Artista.Biografia,
			Foto:            mo.Artista.Foto,
		},
		Genero: &model.Genero{
			ID:     fmt.Sprintf("%d", mo.Genero.ID),
			Nombre: mo.Genero.Nombre,
		},
	}
}

// Trae todas las obras usando paginacion basica de MongoDB
func (r *queryResolver) Obras(ctx context.Context, limit *int32, offset *int32) ([]*model.Obra, error) {
	findOptions := options.Find()
	if offset != nil { findOptions.SetSkip(int64(*offset)) }
	if limit != nil {  findOptions.SetLimit(int64(*limit)) }

	// conectamos con OBRA ULTIMATE HUMONGOSAURIO
	cursor, err := mongodb.GetMongoDB().Collection("obra_ultimate").Find(ctx, bson.D{}, findOptions)
	if err != nil { return nil, err }
	defer cursor.Close(ctx)

	var resultados []models_mongodb.ObraMongo
	if err := cursor.All(ctx, &resultados); err != nil { return nil, err }

	// Itera los resultados usando el helper de conversion para limpiar el codigo
	var obras []*model.Obra
	for _, mo := range resultados {
		obras = append(obras, toGraphQLObra(&mo)) 
	}
	return obras, nil
}

// Busca una sola obra por su ID numerico unico
func (r *queryResolver) FindObra(ctx context.Context, id string) ([]*model.Obra, error) {
	idNumerico, err := strconv.Atoi(id)
	if err != nil { return nil, fmt.Errorf("el ID debe ser valido") }

	var mo models_mongodb.ObraMongo
	err = mongodb.GetMongoDB().Collection("obra_ultimate").FindOne(ctx, bson.D{{Key: "_id", Value: int32(idNumerico)}}).Decode(&mo)
	if err != nil {
		if err == mongo.ErrNoDocuments { return nil, fmt.Errorf("obra no encontrada") }
		return nil, err
	}
	return []*model.Obra{toGraphQLObra(&mo)}, nil
}

// Registra una obra nueva inyectando los datos completos de artista y genero de forma embebida. Es por esto que digo que NO debemos eliminar artista ni genero en Atlas
func (r *mutationResolver) CreateObra(ctx context.Context, input model.NewObra) (*model.Obra, error) {
	if input.Nombre == "" || input.Foto == "" || input.IDArtista == "" || input.IDGenero == "" {
		return nil, fmt.Errorf("hay campos obligatorios vacios")
	}
	if !validStatuses[string(input.Status)] { return nil, fmt.Errorf("status invalido") }

	db := mongodb.GetMongoDB()
	artistaID, _ := strconv.Atoi(input.IDArtista)
	generoID, _ := strconv.Atoi(input.IDGenero)

	// Valida y obtiene los datos del artista original para copiarlos en el documento final
	var artistaData models_mongodb.ArtistaMongo
	if err := db.Collection("artista").FindOne(ctx, bson.D{{Key: "_id", Value: int32(artistaID)}}).Decode(&artistaData); err != nil {
		return nil, fmt.Errorf("el artista no existe")
	}
	
	// Valida y obtiene los datos del genero original para copiarlos en el documento final
	var generoData models_mongodb.GeneroMongo
	if err := db.Collection("genero").FindOne(ctx, bson.D{{Key: "_id", Value: int32(generoID)}}).Decode(&generoData); err != nil {
		return nil, fmt.Errorf("el genero no existe")
	}

	// Calcula el siguiente ID de forma autoincremental buscando el ultimo registro ordenado por ID. NO ME PREGUNTES, ESTE NO SE COMO HACERLO DE OTRA FORMA EN MONGO, ASI QUE AHI VA ESTE HORROR
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	var ultimaObra struct{ ID int32 `bson:"_id"` }
	var nuevoID int32 = 1
	if err := db.Collection("obra_ultimate").FindOne(ctx, bson.D{}, opts).Decode(&ultimaObra); err == nil {
		nuevoID = ultimaObra.ID + 1
	}

	// Crea la estructura completa mapeando los subdocumentos que van embebidos directamente en MongoDB
	nuevaObraDoc := models_mongodb.ObraMongo{
		ID:            nuevoID,
		Nombre:        input.Nombre,
		IDArtista:     int32(artistaID),
		IDGenero:      int32(generoID),
		PrecioObra:    int32(input.Precio),
		FechaCreacion: input.FechaCreacion,
		Status:        string(input.Status),
		Foto:          input.Foto,
		Artista:       artistaData,
		Genero:        generoData,
	}

	if _, err := db.Collection("obra_ultimate").InsertOne(ctx, nuevaObraDoc); err != nil { return nil, err }
	return toGraphQLObra(&nuevaObraDoc), nil
}

// Modifica campos de forma dinamica usando bson.M para armar el mapa de actualizaciones de manera compacta
func (r *mutationResolver) UpdateObra(ctx context.Context, input model.UpdateObra) (*model.Obra, error) {
	idNumerico, err := strconv.Atoi(input.ID)
	if err != nil { return nil, fmt.Errorf("ID invalido") }

	// Se crea un mapa vacio y se agregan propiedades solo si vienen informadas desde el cliente de GraphQL
	updateFields := bson.M{}
	if input.Nombre != nil {  updateFields["nombre"] = *input.Nombre }
	if input.Precio != nil {  updateFields["precio_obra"] = int32(*input.Precio) }
	if input.Foto != nil {    updateFields["foto"] = *input.Foto }
	if input.Status != nil {  updateFields["status"] = string(*input.Status) }
	if input.FechaCreacion != nil && *input.FechaCreacion != "" { updateFields["fecha_creacion"] = *input.FechaCreacion }

	coll := mongodb.GetMongoDB().Collection("obra_ultimate")
	if len(updateFields) > 0 {
		if _, err := coll.UpdateOne(ctx, bson.M{"_id": int32(idNumerico)}, bson.M{"$set": updateFields}); err != nil { return nil, err }
	}

	// Recupera el documento modificado desde la base de datos para retornar la informacion actualizada
	var mo models_mongodb.ObraMongo
	if err := coll.FindOne(ctx, bson.M{"_id": int32(idNumerico)}).Decode(&mo); err != nil { return nil, err }
	return toGraphQLObra(&mo), nil
}

// Elimina un documento por completo usando su clave identificadora unica
func (r *mutationResolver) KillObra(ctx context.Context, id string) (bool, error) {
	idNumerico, err := strconv.Atoi(id)
	if err != nil { return false, fmt.Errorf("ID invalido") }

	res, err := mongodb.GetMongoDB().Collection("obra_ultimate").DeleteOne(ctx, bson.M{"_id": int32(idNumerico)})
	return res.DeletedCount > 0, err
}

func enviarAGraphql(ctx context.Context, limit *int32, offset *int32, pipeline mongo.Pipeline) ([]*model.Obra, error) {
	// si se ingresa un offset o limit
	if offset != nil {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: *offset}})
	}
	if limit != nil {
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: *limit}})
	}

	// conectamos con OBRA ULTIMATE HUMONGOSAURIO
	// y cargamos el pipeline con ordenamiento por precio
	cursor, err := mongodb.GetMongoDB().Collection("obra_ultimate").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// blablabla metemos los resultado a graphql
	var resultados []models_mongodb.ObraMongo
	if err := cursor.All(ctx, &resultados); err != nil {
		return nil, err
	}

	var obras []*model.Obra
	for _, mo := range resultados {
		obras = append(obras, toGraphQLObra(&mo))
	}

	return obras, nil
}
func (r *queryResolver) ObrasPorPrecio(ctx context.Context, limit *int32, offset *int32) ([]*model.Obra, error) {
	// pipeline para ordenar por precio
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$sort", Value: bson.D{{Key: "precio_obra", Value: 1}}}},
	}

	obras, err := enviarAGraphql(ctx, limit, offset, pipeline)
	if err != nil {
		return nil, err
	}

	return obras, nil
}

func (r *queryResolver) ObrasPorPrecioDesc(ctx context.Context, limit *int32, offset *int32) ([]*model.Obra, error) {
	// pipeline para ordenar por precio
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$sort", Value: bson.D{{Key: "precio_obra", Value: -1}}}},
	}

	obras, err := enviarAGraphql(ctx, limit, offset, pipeline)
	if err != nil {
		return nil, err
	}

	return obras, nil
}

func (r *queryResolver) ObrasPorGenero(ctx context.Context, idGenero *int32, limit *int32, offset *int32) ([]*model.Obra, error) {
	// pipeline para filtrar por genero
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "genero._id", Value: idGenero}}}},
	}

	obras, err := enviarAGraphql(ctx, limit, offset, pipeline)
	if err != nil {
		return nil, err
	}

	return obras, nil
}

func (r *queryResolver) ObrasPorDisponibilidad(ctx context.Context, limit *int32, offset *int32) ([]*model.Obra, error) {
	// pipeline para filtrar por status
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "status", Value: "DISPONIBLE"}}}},
	}	

	obras, err := enviarAGraphql(ctx, limit, offset, pipeline)
	if err != nil {
		return nil, err
	}	

	return obras, nil
}


// Funcion auxiliar para descargar una imagen de internet y guardarla localmente
func descargarFotoLocal(url string, nombreArchivo string) error {
	// Asegurar que la ruta de carpetas images/obras exista
	carpeta := filepath.Join("images", "obras")
	if err := os.MkdirAll(carpeta, os.ModePerm); err != nil {
		return err
	}

	// Crear el archivo local vacio
	rutaCompleta := filepath.Join(carpeta, nombreArchivo)
	archivo, err := os.Create(rutaCompleta)
	if err != nil {
		return err
	}
	defer archivo.Close()

	// Hacer la peticion HTTP para bajar la imagen
	respuesta, err := http.Get(url)
	if err != nil {
		return err
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("error al descargar: status %s", respuesta.Status)
	}

	// Copiar el contenido de la respuesta de internet al archivo local
	_, err = io.Copy(archivo, respuesta.Body)
	return err
}

// GenerarObrasMasivas inyecta N cantidad de documentos y descarga sus fotos automaticamente
func (r *mutationResolver) GenerarObrasMasivas(ctx context.Context, cantidad int32) (string, error) {
	if cantidad <= 0 {
		return "", fmt.Errorf("la cantidad debe ser mayor a cero")
	}
	cantidad_int := int(cantidad)
	db := mongodb.GetMongoDB()

	// 1. Traer todos los artistas existentes de la base de datos
	cursorArtistas, err := db.Collection("artista").Find(ctx, bson.D{})
	if err != nil {
		return "", fmt.Errorf("error al buscar artistas: %v", err)
	}
	var artistas []models_mongodb.ArtistaMongo
	if err := cursorArtistas.All(ctx, &artistas); err != nil || len(artistas) == 0 {
		return "", fmt.Errorf("necesitas tener al menos un artista creado en Atlas")
	}

	// 2. Traer todos los generos existentes de la base de datos
	cursorGeneros, err := db.Collection("genero").Find(ctx, bson.D{})
	if err != nil {
		return "", fmt.Errorf("error al buscar generos: %v", err)
	}
	var generos []models_mongodb.GeneroMongo
	if err := cursorGeneros.All(ctx, &generos); err != nil || len(generos) == 0 {
		return "", fmt.Errorf("necesitas tener al menos un genero creado en Atlas")
	}

	// 3. Buscar el ultimo ID correlativo de obra_ultimate
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})
	var ultimaObra struct{ ID int32 `bson:"_id"` }
	var siguienteID int32 = 1
	if err := db.Collection("obra_ultimate").FindOne(ctx, bson.D{}, opts).Decode(&ultimaObra); err == nil {
		siguienteID = ultimaObra.ID + 1
	}

	tiposObra := []string{"Pintura", "Retrato", "Mural", "Boceto", "Fotografia"}
	variantes := []string{"Alpha", "Beta", "Gamma", "Delta", "Omega", "Neon", "Abstracta", "Prime"}

	var nuevosDocumentos []interface{}

	for i := 0; i < cantidad_int; i++ {
		artistaSeleccionado := artistas[i%len(artistas)]
		generoSeleccionado := generos[i%len(generos)]

		tipo := tiposObra[i%len(tiposObra)]
		variante := variantes[(i+1)%len(variantes)]
		
		nombreObra := fmt.Sprintf("%s %s %d", tipo, variante, siguienteID)
		precio := int32(1500 + (i * 35) % 8500) 
		anio := 1850 + (i % 176) 

		// Definir el nombre real del archivo en disco (ej: Pintura101.jpg)
		nombreArchivoFoto := fmt.Sprintf("%s%d.jpg", tipo, siguienteID)

		// URL de una imagen aleatoria de internet (600x400 pixeles)
		// Le agregamos ?random=ID para que internet nos de una foto diferente cada vez
		urlImagenInternet := fmt.Sprintf("https://picsum.photos/600/400?random=%d", siguienteID)

		// Go se encarga de descargarla y guardarla en tu carpeta local automaticamente
		_ = descargarFotoLocal(urlImagenInternet, nombreArchivoFoto)

		doc := models_mongodb.ObraMongo{
			ID:            siguienteID,
			Nombre:        nombreObra,
			IDArtista:     artistaSeleccionado.ID,
			IDGenero:      generoSeleccionado.ID,
			PrecioObra:    precio,
			FechaCreacion: fmt.Sprintf("%d-05-15", anio),
			Status:        "DISPONIBLE",
			Foto:          fmt.Sprintf("images\\obras\\%s", nombreArchivoFoto),
			Artista:       artistaSeleccionado,
			Genero:        generoSeleccionado,
		}

		nuevosDocumentos = append(nuevosDocumentos, doc)
		siguienteID++
	}

	// 5. Insercion masiva a MongoDB Atlas
	resultado, err := db.Collection("obra_ultimate").InsertMany(ctx, nuevosDocumentos)
	if err != nil {
		return "", fmt.Errorf("error en la insercion masiva: %v", err)
	}

	return fmt.Sprintf("Se insertaron %d obras y se descargaron sus imagenes", len(resultado.InsertedIDs)), nil
}

func (r *mutationResolver) SincronizarFotosImgBBMasivo(ctx context.Context) (string, error) {
	db := mongodb.GetMongoDB()
	collObra := db.Collection("obra_ultimate")

	// 1. Leer todas las líneas de enlaces.txt
	contenido, err := os.ReadFile("enlaces.txt")
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el archivo enlaces.txt: %v", err)
	}

	// 2. Traer todas las obras para mapear el estado actual en Atlas
	cursorTodos, err := collObra.Find(ctx, bson.M{})
	if err != nil {
		return "", fmt.Errorf("error al mapear enlaces existentes: %v", err)
	}
	defer cursorTodos.Close(ctx)

	enlacesUsados := make(map[string]bool)
	var obrasPendientes []bson.M

	for cursorTodos.Next(ctx) {
		var obra bson.M
		if err := cursorTodos.Decode(&obra); err != nil {
			continue
		}

		fotoStr := fmt.Sprintf("%v", obra["foto"])
		if strings.HasPrefix(strings.ToLower(fotoStr), "http") {
			enlacesUsados[strings.TrimSpace(fotoStr)] = true
		} else if strings.HasPrefix(strings.ToLower(fotoStr), "images") {
			obrasPendientes = append(obrasPendientes, obra)
		}
	}

	// Compilar regex para buscar EXACTAMENTE 3 numeros seguidos antes de la extensión de la imagen
	reTresCifras := regexp.MustCompile(`[0-9]{3}\.[a-zA-Z]+$`)

	// 3. Filtrar el archivo TXT para separar los enlaces NUEVOS y el grupo de RECICLAJE (3 cifras)
	var enlacesNuevosDisponibles []string
	var poolReciclajeTresCifras []string

	for _, l := range strings.Split(string(contenido), "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		// Si cumple con tener un numero de 3 cifras al final, entra al pool de reciclaje
		if reTresCifras.MatchString(l) {
			poolReciclajeTresCifras = append(poolReciclajeTresCifras, l)
		}

		// Si ademas no ha sido usado, se usa como cartucho prioritario
		if !enlacesUsados[l] {
			enlacesNuevosDisponibles = append(enlacesNuevosDisponibles, l)
		}
	}

	log.Printf(" INICIANDO MIGRACIÓN CON RECICLAJE FILTRADO ")
	log.Printf("Obras por arreglar: %d", len(obrasPendientes))
	log.Printf("Enlaces nuevos disponibles: %d", len(enlacesNuevosDisponibles))
	log.Printf("Enlaces elegibles para reciclaje (3 cifras): %d", len(poolReciclajeTresCifras))

	if len(poolReciclajeTresCifras) == 0 {
		return "", fmt.Errorf("error: no se encontraron enlaces con nombres de 3 cifras al final en enlaces.txt")
	}

	actualizadas := 0

	// 4. Inyectar enlaces priorizando los nuevos, y reciclando los de 3 cifras si se acaban
	for i, obra := range obrasPendientes {
		idObra := obra["_id"]
		var urlNubeAsignada string

		if i < len(enlacesNuevosDisponibles) {
			// Usamos los ultimos cartuchos limpios que queden
			urlNubeAsignada = enlacesNuevosDisponibles[i]
		} else {
			// Bucle circular sobre el pool de 3 cifras usando el operador residuo (%)
			indiceReciclado := (i - len(enlacesNuevosDisponibles)) % len(poolReciclajeTresCifras)
			urlNubeAsignada = poolReciclajeTresCifras[indiceReciclado]
		}

		_, errUpdate := collObra.UpdateOne(ctx, 
			bson.M{"_id": idObra}, 
			bson.M{"$set": bson.M{"foto": urlNubeAsignada}},
		)
		if errUpdate == nil {
			actualizadas++
		}
	}

	log.Printf(" FIN: Se forzo la nube en las %d obras rezagadas ", actualizadas)
	return fmt.Sprintf("¡Operacion exitosa! Se pasaron %d obras a la nube reciclando las fotos de 3 cifras.", actualizadas), nil
}