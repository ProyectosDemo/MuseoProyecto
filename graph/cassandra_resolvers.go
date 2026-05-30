package graph

import (
	"context"
	"fmt"
	"main/data_bases/cassandra"
	"main/graph/model"
	"time"

	"github.com/gocql/gocql"
)

// EmitirFacturaHistorica guarda un registro financiero desnormalizado en Cassandra
func (r *mutationResolver) EmitirFacturaHistorica(ctx context.Context, periodo string, fechaFactura string, idOrden string, idCliente string, clienteNombre string, idObra string, obraNombre string, precio float64) (*model.FacturaHistorica, error) {
	
	// 1. Parsear la fecha del Frontend
	fechaParsed, err := time.Parse("2006-01-02", fechaFactura)
	if err != nil {
		fechaParsed, err = time.Parse(time.RFC3339, fechaFactura)
		if err != nil {
			fechaParsed = time.Now() 
		}
	}

	// 2. Reglas de negocio calculadas en el Backend
	const porcentajeIVA = 0.21   // Actualizado al 21% que usas en el Frontend
	const porcentajeMuseo = 0.10 
	
	iva := precio * porcentajeIVA
	gananciaMuseo := precio * porcentajeMuseo
	total := precio + iva

	session := cassandra.GetCassandra()

	// 3. Query de insercion (Ya no usamos ParseUUID, enviamos el string directo)
	query := `INSERT INTO facturas_por_periodo (
		periodo, fecha_factura, id_orden, id_cliente, cliente_nombre, 
		id_obra, obra_nombre, precio, iva, total, ganancia_museo
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	err = session.Query(query, 
		periodo, fechaParsed, idOrden, idCliente, clienteNombre, 
		idObra, obraNombre, precio, iva, total, gananciaMuseo,
	).Exec()

	if err != nil {
		return nil, fmt.Errorf("error al persistir factura en Cassandra: %v", err)
	}

	return &model.FacturaHistorica{
		Periodo:        periodo,
		FechaFactura:   fechaParsed.Format("2006-01-02"),
		IDOrden:        idOrden,
		IDCliente:      idCliente,
		ClienteNombre:  clienteNombre,
		IDObra:         idObra,
		ObraNombre:     obraNombre,
		Precio:         precio,
		Iva:            iva,
		Total:          total,
		GananciaMuseo:  gananciaMuseo,
	}, nil
}

// ObtenerReporteFacturas extrae de forma masiva los datos agrupados por periodo
func (r *queryResolver) ObtenerReporteFacturas(ctx context.Context, periodo string) ([]*model.FacturaHistorica, error) {
	session := cassandra.GetCassandra()
	var reportes []*model.FacturaHistorica

	query := `SELECT id_orden, fecha_factura, id_cliente, cliente_nombre, id_obra, obra_nombre, precio, iva, total, ganancia_museo 
	          FROM facturas_por_periodo WHERE periodo = ?`
	
	iter := session.Query(query, periodo).Iter()
	
	// Escaneamos los IDs foraneos directamente como strings
	var idOrden, idCliente, idObra string
	var fechaFactura time.Time
	var clienteNombre, obraNombre string
	var precio, iva, total, gananciaMuseo float64

	for iter.Scan(&idOrden, &fechaFactura, &idCliente, &clienteNombre, &idObra, &obraNombre, &precio, &iva, &total, &gananciaMuseo) {
		reportes = append(reportes, &model.FacturaHistorica{
			Periodo:        periodo,
			FechaFactura:   fechaFactura.Format("2006-01-02"),
			IDOrden:        idOrden,
			IDCliente:      idCliente,
			ClienteNombre:  clienteNombre,
			IDObra:         idObra,
			ObraNombre:     obraNombre,
			Precio:         precio,
			Iva:            iva,
			Total:          total,
			GananciaMuseo:  gananciaMuseo,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error leyendo datos de Cassandra: %v", err)
	}

	return reportes, nil
}

// RegistrarEventoObra inserta un registro inmutable en la bitacora
func (r *mutationResolver) RegistrarEventoObra(ctx context.Context, input model.NewBitacoraObra) (*model.BitacoraObra, error) {
	
	// El evento local de Cassandra SI necesita ser TimeUUID para evitar colisiones
	idEvento := gocql.TimeUUID()
	fechaEvento := time.Now()

	session := cassandra.GetCassandra()

	query := `INSERT INTO bitacora_obras (
		id_obra, fecha_evento, id_evento, tipo_evento, 
		descripcion, usuario_responsable, estatus_anterior, estatus_nuevo
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	// Enviamos input.IDObra directo sin parsearlo a UUID
	err := session.Query(query,
		input.IDObra, fechaEvento, idEvento, input.TipoEvento,
		input.Descripcion, input.UsuarioResponsable, input.EstatusAnterior, input.EstatusNuevo,
	).Exec()

	if err != nil {
		return nil, fmt.Errorf("error al registrar evento en Cassandra: %v", err)
	}

	return &model.BitacoraObra{
		IDObra:             input.IDObra,
		FechaEvento:        fechaEvento.Format(time.RFC3339),
		IDEvento:           idEvento.String(),
		TipoEvento:         input.TipoEvento,
		Descripcion:        input.Descripcion,
		UsuarioResponsable: input.UsuarioResponsable,
		EstatusAnterior:    input.EstatusAnterior,
		EstatusNuevo:       input.EstatusNuevo,
	}, nil
}

// ObtenerBitacoraObra lee el historial de auditoria
func (r *queryResolver) ObtenerBitacoraObra(ctx context.Context, idObra string) ([]*model.BitacoraObra, error) {
	session := cassandra.GetCassandra()
	var eventos []*model.BitacoraObra

	query := `SELECT id_evento, fecha_evento, tipo_evento, descripcion, usuario_responsable, estatus_anterior, estatus_nuevo 
	          FROM bitacora_obras WHERE id_obra = ?`

	iter := session.Query(query, idObra).Iter()

	var idEvento gocql.UUID
	var fechaEvento time.Time
	var tipoEvento, descripcion, usuarioResponsable, estatusAnterior, estatusNuevo string

	for iter.Scan(&idEvento, &fechaEvento, &tipoEvento, &descripcion, &usuarioResponsable, &estatusAnterior, &estatusNuevo) {
		eventos = append(eventos, &model.BitacoraObra{
			IDObra:             idObra,
			FechaEvento:        fechaEvento.Format(time.RFC3339),
			IDEvento:           idEvento.String(),
			TipoEvento:         tipoEvento,
			Descripcion:        descripcion,
			UsuarioResponsable: usuarioResponsable,
			EstatusAnterior:    estatusAnterior,
			EstatusNuevo:       estatusNuevo,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("error leyendo bitacora de Cassandra: %v", err)
	}

	return eventos, nil
}