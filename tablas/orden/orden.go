package orden

import (
	"main/middleware"
	"main/models"
	"main/mysql"
	"main/tablas/cliente"
	"main/tablas/obra"
	"main/tablas/trabajador"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
)

var EstatusOrdenEnum = graphql.NewEnum(graphql.EnumConfig{
	Name: "EstatusOrden",
	Values: graphql.EnumValueConfigMap{
		"Pendiente": &graphql.EnumValueConfig{
			Value: models.Pendiente,
		},
		"Concretada": &graphql.EnumValueConfig{
			Value: models.Pagado,
		},
		"Cancelada": &graphql.EnumValueConfig{
			Value: models.Cancelado,
		},
	},
})

func CreateOrdenType(clienteType *graphql.Object, obraType *graphql.Object, trabajadorType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Orden",
		Fields: graphql.Fields{
			"id_orden":      &graphql.Field{Type: graphql.Int},
			"id_cliente":    &graphql.Field{Type: graphql.Int},
			"id_obra":       &graphql.Field{Type: graphql.Int},
			"id_trabajador": &graphql.Field{Type: graphql.Int},

			// esto se calcula en workbench automaticamente 
			"precio_obra": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if orden, ok := p.Source.(models.Orden); ok {
						return orden.Precio_obra, nil
					}
					return nil, nil
				},
			},
			"iva": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if orden, ok := p.Source.(models.Orden); ok {
						return orden.Iva, nil
					}
					return nil, nil
				},
			},
			"ganancia_museo": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if orden, ok := p.Source.(models.Orden); ok {
						return orden.Ganancia_museo, nil
					}
					return nil, nil
				},
			},
			"total": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if orden, ok := p.Source.(models.Orden); ok {
						return orden.Total, nil
					}
					return nil, nil
				},
			},

			"fecha_orden": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if orden, ok := p.Source.(models.Orden); ok {
						return orden.Fecha_orden.Format("02-01-2006"), nil
					}
					return nil, nil
				},
			},

			"estatus": &graphql.Field{Type: EstatusOrdenEnum},

			"cliente": &graphql.Field{
				Type: clienteType,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if orden, ok := p.Source.(models.Orden); ok {
						return cliente.GetClienteByID(int(orden.Id_cliente))
					}
					return nil, nil
				},
			},
			"obra": &graphql.Field{
				Type: obraType,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if orden, ok := p.Source.(models.Orden); ok {
						return obra.GetObraByID(int(orden.Id_obra))
					}
					return nil, nil
				},
			},
			"trabajador": &graphql.Field{
				Type: trabajadorType,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					if orden, ok := p.Source.(models.Orden); ok {
						return trabajador.GetTrabajadorByID(int(orden.Id_trabajador))
					}
					return nil, nil
				},
			},
		},
	})
}

func CreateOrdenField(ordenType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        ordenType,
		Description: "Crear una nueva orden",
		Args: graphql.FieldConfigArgument{
			"id_cliente":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"id_obra":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"id_trabajador": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"fecha_orden":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			"estatus":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(EstatusOrdenEnum)},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id_cliente := p.Args["id_cliente"].(int)
			id_obra := p.Args["id_obra"].(int)
			id_trabajador := p.Args["id_trabajador"].(int)
			fecha_orden_str := p.Args["fecha_orden"].(string)
			estatus := p.Args["estatus"].(models.EstatusOrden)

			// Validaciones
			errores := []string{}
			if !cliente.ClienteExiste(id_cliente) {
				errores = append(errores, "Error: El cliente con ID "+strconv.Itoa(id_cliente)+" no existe.")
			}
			if !obra.ObraExiste(id_obra) {
				errores = append(errores, "Error: La obra con ID "+strconv.Itoa(id_obra)+" no existe.")
			}
			if !trabajador.TrabajadorExiste(id_trabajador) {
				errores = append(errores, "Error: El trabajador con ID "+strconv.Itoa(id_trabajador)+" no existe.")
			}
			if len(errores) != 0 {
				for _, e := range errores {
					println(e)
				}
				return nil, nil
			}

			fecha_orden, err := time.Parse("2006-01-02", fecha_orden_str)
			middleware.PanicButton(err)

			// Obtener precio_obra desde la obra
			var precio_obra float64
			err = mysql.GetBD().QueryRow("SELECT precio_obra FROM obra WHERE id_obra = ?", id_obra).Scan(&precio_obra)
			middleware.PanicButton(err)

			// Calcular ganancia del museo
			ganancia_museo := precio_obra * 0.2 // ejemplo: 20%

			// Insertar directamente en la DB (iva y total son generados por MySQL)
			id := mysql.Insertar(
				mysql.GetBD(),
				`INSERT INTO orden 
					(id_cliente, id_obra, id_trabajador, precio_obra, ganancia_museo, fecha_orden, estatus)
					VALUES (?, ?, ?, ?, ?, ?, ?)`,
				id_cliente, id_obra, id_trabajador, precio_obra, ganancia_museo, fecha_orden, estatus,
			)

			return models.Orden{
				Id_orden:       id,
				Id_cliente:     int64(id_cliente),
				Id_obra:        int64(id_obra),
				Id_trabajador:  int64(id_trabajador),
				Precio_obra:         precio_obra,
				Ganancia_museo: ganancia_museo,
				Fecha_orden:    fecha_orden,
				Estatus:        estatus,
			}, nil
		},
	}
}

func GetOrdenesField(ordenType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(ordenType),
		Description: "Lista de ordenes",
		Args: graphql.FieldConfigArgument{
			"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
			"offset": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			limit, _ := p.Args["limit"].(int)
			if limit <= 0 || limit > 20 {
				limit = 10
			}
			offset, _ := p.Args["offset"].(int)
			if offset < 0 {
				offset = 0
			}
			return GetOrdenes(limit, offset)
		},
	}
}


func GetOrdenes(limit int, offset int) ([]models.Orden, error) {
	var ordenes []models.Orden
	rows, err := mysql.GetBD().Query(
    `SELECT 
		id_orden, id_cliente, id_obra, id_trabajador, precio_obra, iva, ganancia_museo, total, fecha_orden, estatus
		FROM orden
		ORDER BY id_orden
		LIMIT ` + strconv.Itoa(limit) + ` OFFSET ` + strconv.Itoa(offset),
	)
	middleware.PanicButton(err)
	defer rows.Close()

	for rows.Next() {
		var o models.Orden
		if err := rows.Scan(
			&o.Id_orden, &o.Id_cliente, &o.Id_obra, &o.Id_trabajador, &o.Precio_obra, &o.Iva, &o.Ganancia_museo, &o.Total, &o.Fecha_orden, &o.Estatus,
		); err != nil {
			return nil, err
		}
		ordenes = append(ordenes, o)
	}
	return ordenes, nil
}

func UpdateOrdenField(ordenType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        ordenType,
		Description: "Actualizar una orden por ID",
		Args: graphql.FieldConfigArgument{
			"id_orden":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"id_cliente":    &graphql.ArgumentConfig{Type: graphql.Int},
			"id_obra":       &graphql.ArgumentConfig{Type: graphql.Int},
			"id_trabajador": &graphql.ArgumentConfig{Type: graphql.Int},
			"ganancia_museo": &graphql.ArgumentConfig{Type: graphql.Float},
			"fecha_orden":   &graphql.ArgumentConfig{Type: graphql.String},
			"estatus": &graphql.ArgumentConfig{Type: EstatusOrdenEnum},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id_orden := p.Args["id_orden"].(int)

			var o models.Orden
			err := mysql.GetBD().QueryRow(
				`SELECT id_orden, id_cliente, id_obra, id_trabajador, precio_obra, ganancia_museo, fecha_orden, estatus
				 FROM orden WHERE id_orden = ?`,
				id_orden,
			).Scan(&o.Id_orden, &o.Id_cliente, &o.Id_obra, &o.Id_trabajador, &o.Precio_obra, &o.Ganancia_museo, &o.Fecha_orden, &o.Estatus)
			middleware.PanicButton(err)

			if id_cliente, ok := p.Args["id_cliente"].(int); ok {
				o.Id_cliente = int64(id_cliente)
			}
			if id_obra, ok := p.Args["id_obra"].(int); ok {
				o.Id_obra = int64(id_obra)
				err := mysql.GetBD().QueryRow("SELECT precio_obra FROM obra WHERE id_obra = ?", id_obra).Scan(&o.Precio_obra)
				middleware.PanicButton(err)
			}
			if id_trabajador, ok := p.Args["id_trabajador"].(int); ok {
				o.Id_trabajador = int64(id_trabajador)
			}
			if ganancia_museo, ok := p.Args["ganancia_museo"].(float64); ok {
				o.Ganancia_museo = ganancia_museo
			}
			if fechaStr, ok := p.Args["fecha_orden"].(string); ok && fechaStr != "" {
				fechaParseada, err := time.Parse("2006-01-02", fechaStr)
				middleware.PanicButton(err)
				o.Fecha_orden = fechaParseada
			}
			if estatusVal, ok := p.Args["estatus"].(string); ok && estatusVal != "" {
				o.Estatus = models.EstatusOrden(estatusVal)
			}

			_, err = mysql.GetBD().Exec(
				`UPDATE orden SET 
					id_cliente=?, id_obra=?, id_trabajador=?, precio_obra=?, ganancia_museo=?, fecha_orden=?, estatus=? 
				 WHERE id_orden=?`,
				o.Id_cliente, o.Id_obra, o.Id_trabajador, o.Precio_obra, o.Ganancia_museo, o.Fecha_orden, o.Estatus, id_orden,
			)
			middleware.PanicButton(err)

			return o, nil
		},
	}
}


func DeleteOrdenField(ordenType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        ordenType,
		Description: "Eliminar una orden por ID",
		Args: graphql.FieldConfigArgument{
			"id_orden": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(p graphql.ResolveParams) (any, error) {
			id_orden := p.Args["id_orden"].(int)

			var o models.Orden
			err := mysql.GetBD().QueryRow(
				`SELECT id_orden, id_cliente, id_obra, id_trabajador, precio_obra, iva, ganancia_museo, total, fecha_orden, estatus
				 FROM orden WHERE id_orden=?`,
				id_orden,
			).Scan(&o.Id_orden, &o.Id_cliente, &o.Id_obra, &o.Id_trabajador, &o.Precio_obra, &o.Iva, &o.Ganancia_museo, &o.Total, &o.Fecha_orden,
				&o.Estatus)
			middleware.PanicButton(err)

			_, err = mysql.GetBD().Exec(
				"DELETE FROM orden WHERE id_orden=?",
				id_orden,
			)
			middleware.PanicButton(err)

			return o, nil
		},
	}
}
