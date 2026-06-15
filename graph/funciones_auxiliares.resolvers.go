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
	"strings"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


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


func (r *queryResolver) ObtenerRecomendaciones(ctx context.Context, idCliente string) ([]*model.Obra, error) {
	session := r.Resolver.Neo4j.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "museoproyecto",
		AccessMode:   neo4j.AccessModeRead,
	})
	defer session.Close(ctx)

	query := `
		MATCH (c:Comprador {id: $idCliente})-[:COMPRÓ]->(:Obra)
		MATCH (c)-[:SUGERIDA_POR_GENERO]->(o:Obra)
		WHERE o.status = "DISPONIBLE"
		RETURN DISTINCT o.id AS id, o.nombre AS nombre, o.foto AS foto
		LIMIT 10
	`

	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		var idClienteInterface interface{}
		var idInt int
		if _, err := fmt.Sscanf(idCliente, "%d", &idInt); err == nil {
			idClienteInterface = idInt
		} else {
			idClienteInterface = idCliente
		}

		result, txErr := tx.Run(ctx, query, map[string]interface{}{"idCliente": idClienteInterface})
		if txErr != nil {
			return nil, txErr
		}

		var recomendaciones []*model.Obra
		for result.Next(ctx) {
			record := result.Record()
			id, _ := record.Get("id")
			nombre, _ := record.Get("nombre")
			foto, _ := record.Get("foto")

			obra := &model.Obra{
				ID:     fmt.Sprintf("%v", id), 
				Nombre: fmt.Sprintf("%v", nombre),
				Foto:   fmt.Sprintf("%v", foto),
			}
			recomendaciones = append(recomendaciones, obra)
		}
		return recomendaciones, nil
	})

	if err != nil {
		return nil, err
	}

	return res.([]*model.Obra), nil
}