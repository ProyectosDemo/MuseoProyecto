document.addEventListener('DOMContentLoaded', cargarOrdenesConcretadas);

const navMenu = document.getElementById("navMenu");
const trabajadorId = localStorage.getItem("trabajadorId");
const trabajadorNombre = localStorage.getItem("trabajadorNombre");

if (!trabajadorId) location.href = "login.html";

const logoutBtn = document.createElement("a");
logoutBtn.href = "#";
logoutBtn.textContent = "Cerrar Sesión";
logoutBtn.style.marginLeft = "10px";

logoutBtn.addEventListener("click", () => {
    localStorage.removeItem("trabajadorNombre");
    localStorage.removeItem("trabajadorId");
    localStorage.removeItem("trabajadorAdmin");
    location.reload();
});
navMenu.appendChild(logoutBtn);

document.getElementById("menuOrdenes").addEventListener("click", () => mostrarPanel("ordenesPanel"));
document.getElementById("menuClientes").addEventListener("click", () => mostrarPanel("clientesPanel"));
document.getElementById("menuObras").addEventListener("click", () => mostrarPanel("obrasPanel"));
document.getElementById("menuFacturas").addEventListener("click", () => mostrarPanel("facturasPanel"));

function mostrarPanel(panelId) {
    const panels = document.querySelectorAll(".main-content .panel");
    panels.forEach(p => p.hidden = true);
    document.getElementById(panelId).hidden = false;
}

// Seccion de Ordenes Concretadas y guardado en la base de datos Cassandra
async function cargarOrdenesConcretadas() {
    try {
        // Pedimos todas las ordenes al backend de Go
        const query = `query { Ordenes(limit: 100, offset: 0) { id id_obra id_cliente id_trabajador fecha status } }`;
        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query })
        });
        const data = await res.json();
        
        // Filtramos en memoria para quedarnos solo con las que tienen status CONCRETADA
        const ordenes = data?.data?.Ordenes.filter(o => o.status === "CONCRETADA") || [];
        const container = document.getElementById("ordenesContainer");
        container.innerHTML = "";

        // Iteramos sobre cada orden concretada para armar su tarjeta visual
        for (const orden of ordenes) {
            
            // Subconsulta para traer el nombre y precio de la obra desde MySQL/MongoDB
            const obraQuery = `query { findObra(id: "${orden.id_obra}") { id nombre precio } }`;
            const obraRes = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: obraQuery })
            });
            const obraData = await obraRes.json();
            const obra = obraData?.data?.findObra[0];

            // Subconsulta para traer el nombre del cliente
            const clienteQuery = `query { findCliente(id: "${orden.id_cliente}") { id nombre } }`;
            const clienteRes = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: clienteQuery })
            });
            const clienteData = await clienteRes.json();
            const cliente = clienteData?.data?.findCliente[0];

            // Contenedor principal que agrupa la orden y el espacio de la futura factura
            const wrapper = document.createElement("div");
            wrapper.className = "orden-factura-wrapper";

            // Tarjeta que muestra los datos traidos de MySQL
            const ordenCard = document.createElement("div");
            ordenCard.className = "orden-card";
            ordenCard.innerHTML = `
                <p><strong>ID Orden:</strong> ${orden.id}</p>
                <p><strong>Obra:</strong> ${obra?.nombre || "N/D"}</p>
                <p><strong>Cliente:</strong> ${cliente?.nombre || "N/D"}</p>
                <p><strong>Fecha:</strong> ${orden.fecha.split('T')[0]}</p>
                <p><strong>Precio:</strong> $${obra?.precio || 0}</p>
                <button class="facturaBtn">Generar Factura</button>
            `;

            // Espacio vacio donde se dibujara la factura al hacer clic
            const facturaCard = document.createElement("div");
            facturaCard.className = "factura-card";
            facturaCard.innerHTML = `<p>Factura aparecerá aquí al seleccionar la orden.</p>`;

            wrapper.appendChild(ordenCard);
            wrapper.appendChild(facturaCard);
            container.appendChild(wrapper);

            // Logica del boton Generar Factura
            const btnFactura = ordenCard.querySelector(".facturaBtn");
            btnFactura.addEventListener("click", async () => {
                
                // 1. Dibuja la factura en el HTML para que el usuario la vea
                generarFacturaEnCard(orden, obra, cliente, facturaCard);

                // 2. Extraemos el mes y ano (ej. 2026-05) para usarlo como Partition Key en Cassandra
                const periodoFormateado = orden.fecha.substring(0, 7);
                // Extraemos la fecha sin la hora para simplificar el almacenamiento
                const fechaLimpia = orden.fecha.split('T')[0];

                // 3. Preparamos la mutacion de GraphQL para guardar el registro historico
                const mutation = `
                    mutation {
                        emitirFacturaHistorica(
                            periodo: "${periodoFormateado}",
                            fecha_factura: "${fechaLimpia}",
                            id_orden: "${orden.id}",
                            id_cliente: "${orden.id_cliente}",
                            cliente_nombre: "${cliente?.nombre || 'Anónimo'}",
                            id_obra: "${orden.id_obra}",
                            obra_nombre: "${obra?.nombre || 'Obra Desconocida'}",
                            precio: ${obra?.precio || 0}
                        ) {
                            id_orden
                        }
                    }
                `;

                try {
                    // Enviamos la peticion a Go para que la inserte en Cassandra
                    const response = await fetch("http://localhost:8080/query", {
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({ query: mutation })
                    });
                    const resData = await response.json();
                    
                    if (resData.errors) {
                        console.error("Error al guardar en Cassandra:", resData.errors);
                        alert("Error al emitir factura en la base de datos.");
                    } else {
                        console.log("Factura guardada exitosamente en Cassandra", resData.data);
                        // Cambiamos el estilo del boton para evitar clics duplicados
                        btnFactura.textContent = "Factura Emitida";
                        btnFactura.disabled = true;
                        btnFactura.style.backgroundColor = "#ccc";
                        btnFactura.style.cursor = "not-allowed";
                    }
                } catch (error) {
                    console.error("Error de conexion con el backend:", error);
                }
            });
        }
    } catch (err) {
        console.error("Error cargando las ordenes:", err);
    }
}

// Funcion auxiliar que calcula los impuestos y ganancias, y actualiza el HTML de la tarjeta
function generarFacturaEnCard(orden, obra, cliente, card) {
    const IVA = 0.21;
    const porcentajeMuseo = 0.10;
    
    // Calculo de reglas de negocio
    const precio = obra.precio;
    const iva = precio * IVA;
    const gananciaMuseo = precio * porcentajeMuseo;
    const total = precio + iva;

    // Pintamos los datos financieros calculados
    card.innerHTML = `
        <p><strong>Factura ID:</strong> ${orden.id}</p>
        <p><strong>Cliente:</strong> ${cliente.nombre}</p>
        <p><strong>Obra:</strong> ${obra.nombre}</p>
        <p><strong>Precio:</strong> $${precio.toFixed(2)}</p>
        <p><strong>IVA (21%):</strong> $${iva.toFixed(2)}</p>
        <p><strong>Total:</strong> $${total.toFixed(2)}</p>
        <p><strong>Ganancia Museo:</strong> $${gananciaMuseo.toFixed(2)}</p>
        <p><strong>Fecha:</strong> ${new Date(orden.fecha).toISOString().split('T')[0]}</p>
    `;
}

// Seccion de Reportes de Membresias
// Escucha el clic del boton filtrar y valida que haya fechas
document.getElementById("btnFiltrarMembresias").addEventListener("click", async () => {
    const desde = document.getElementById("fechaInicio").value;
    const hasta = document.getElementById("fechaFin").value;
    if (!desde || !hasta) return alert("Por favor selecciona ambas fechas.");
    cargarResumenMembresias(desde, hasta);
});

// Pide las membresias y las filtra segun las fechas elegidas en los calendarios
async function cargarResumenMembresias(fechaInicio, fechaFin) {
    const query = `query { Membresias(limit: 1000, offset: 0) { id_membresia fecha cliente { nombre } tarjeta { numero_tarjeta } } }`;
    const res = await fetch("http://localhost:8080/query", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query })
    });
    const data = await res.json();
    const membresias = data?.data?.Membresias || [];
    
    // Comparamos fechas convirtiendo los strings a objetos Date de Javascript
    const filtradas = membresias.filter(m => {
        const fechaM = new Date(m.fecha);
        return fechaM >= new Date(fechaInicio) && fechaM <= new Date(fechaFin);
    });

    const container = document.getElementById("resumenMembresias");
    container.innerHTML = "<h3>Resumen de Membresías</h3>";
    
    // Si no hay datos, mostramos mensaje y salimos de la funcion
    if (!filtradas.length) return container.innerHTML += "<p>No hay membresías en este periodo.</p>";

    // Generamos una tabla HTML iterando sobre los resultados
    const table = document.createElement("table");
    table.innerHTML = "<tr><th>Cliente</th><th>Fecha de Membresía</th></tr>";
    filtradas.forEach(m => {
        const row = document.createElement("tr");
        row.innerHTML = `<td>${m.cliente.nombre}</td><td>${new Date(m.fecha).toISOString().split('T')[0]}</td>`;
        table.appendChild(row);
    });
    container.appendChild(table);
}

// Seccion de Reportes de Obras Vendidas
// Escucha el clic y valida las fechas
document.getElementById("btnFiltrarObras").addEventListener("click", async () => {
    const desde = document.getElementById("fechaObrasInicio").value;
    const hasta = document.getElementById("fechaObrasFin").value;
    if (!desde || !hasta) return alert("Por favor selecciona ambas fechas.");
    cargarObrasVendidas(desde, hasta);
});

// Calcula y agrupa cuantas veces se vendio una misma obra y cuanto dinero genero
async function cargarObrasVendidas(fechaInicio, fechaFin) {
    const query = `query { Ordenes(limit: 1000, offset: 0) { id id_obra fecha status } }`;
    const res = await fetch("http://localhost:8080/query", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query })
    });
    const data = await res.json();
    
    // Solo nos interesan las ordenes concretadas
    const ordenes = data?.data?.Ordenes.filter(o => o.status === "CONCRETADA") || [];
    
    // Usamos un objeto como mapa para agrupar las obras por nombre
    const obrasVendidasMap = {};

    for (const orden of ordenes) {
        // Descartamos las ordenes que esten fuera del rango de fecha
        const fechaOrden = new Date(orden.fecha);
        if (fechaOrden < new Date(fechaInicio) || fechaOrden > new Date(fechaFin)) continue;

        // Buscamos los detalles de la obra asociada a esta orden
        const obraQuery = `query { findObra(id: "${orden.id_obra}") { nombre precio } }`;
        const obraRes = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query: obraQuery })
        });
        const obraData = await obraRes.json();
        const obra = obraData?.data?.findObra[0];
        
        if (!obra) continue;

        // Si es la primera vez que vemos esta obra, inicializamos sus contadores
        if (!obrasVendidasMap[obra.nombre]) {
            obrasVendidasMap[obra.nombre] = { cantidad: 0, total: 0, precio: obra.precio };
        }
        
        // Sumamos 1 a la cantidad vendida y acumulamos el dinero
        obrasVendidasMap[obra.nombre].cantidad++;
        obrasVendidasMap[obra.nombre].total += obra.precio;
    }

    const container = document.getElementById("obrasContainer");
    container.innerHTML = "<h3>Obras Vendidas por Fecha</h3>";
    const table = document.createElement("table");
    table.innerHTML = "<tr><th>Obra</th><th>Cantidad Vendida</th><th>Precio Unitario</th><th>Total Vendido</th></tr>";

    // Convertimos el mapa en un arreglo y dibujamos cada fila en la tabla
    for (const [nombre, info] of Object.entries(obrasVendidasMap)) {
        const row = document.createElement("tr");
        row.innerHTML = `<td>${nombre}</td><td>${info.cantidad}</td><td>$${info.precio.toFixed(2)}</td><td>$${info.total.toFixed(2)}</td>`;
        table.appendChild(row);
    }
    container.appendChild(table);
}

// Seccion de Reporte de Facturas Masivo utilizando Cassandra
// Escucha el clic y valida las fechas
document.getElementById("btnFiltrarFacturas").addEventListener("click", async () => {
    const desde = document.getElementById("fechaFacturasInicio").value;
    const hasta = document.getElementById("fechaFacturasFin").value;
    if (!desde || !hasta) return alert("Por favor selecciona ambas fechas.");
    cargarResumenFacturas(desde, hasta);
});

// Pide datos a Cassandra de forma super rapida porque la data ya viene plana y desnormalizada
async function cargarResumenFacturas(fechaInicio, fechaFin) {
    
    // Cortamos la fecha "Desde" para obtener el formato YYYY-MM que Cassandra usa como llave de particion
    const periodo = fechaInicio.substring(0, 7); 

    // A diferencia de MySQL, aqui no necesitamos subconsultas, un solo request trae todo
    const query = `
        query {
            obtenerReporteFacturas(periodo: "${periodo}") {
                id_orden
                obra_nombre
                cliente_nombre
                precio
                iva
                total
                ganancia_museo
                fecha_factura
            }
        }
    `;

    try {
        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query })
        });
        const data = await res.json();
        
        const facturasCassandra = data?.data?.obtenerReporteFacturas || [];

        // Filtramos en Javascript para que los dias coincidan exactamente con los calendarios
        const filtradas = facturasCassandra.filter(f => {
            return f.fecha_factura >= fechaInicio && f.fecha_factura <= fechaFin;
        });

        const container = document.getElementById("facturasContainer");
        container.innerHTML = "<h3>Resumen de Facturas (Cargado de Cassandra)</h3>";
        
        if (!filtradas.length) {
            container.innerHTML += "<p>No hay facturas inmutables emitidas en este periodo.</p>";
            return;
        }

        // Variables para sumar totales globales en la parte inferior de la tabla
        let totalFacturas = 0, totalIVA = 0, totalMuseo = 0;

        const table = document.createElement("table");
        table.innerHTML = "<tr><th>ID Orden</th><th>Obra</th><th>Cliente</th><th>Precio</th><th>IVA</th><th>Total</th><th>Ganancia Museo</th><th>Fecha</th></tr>";

        // Iteramos las facturas historicas y acumulamos el dinero global
        filtradas.forEach(f => {
            totalFacturas += f.total;
            totalIVA += f.iva;
            totalMuseo += f.ganancia_museo;

            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${f.id_orden}</td>
                <td>${f.obra_nombre}</td>
                <td>${f.cliente_nombre}</td>
                <td>$${f.precio.toFixed(2)}</td>
                <td>$${f.iva.toFixed(2)}</td>
                <td><strong>$${f.total.toFixed(2)}</strong></td>
                <td>$${f.ganancia_museo.toFixed(2)}</td>
                <td>${f.fecha_factura}</td>
            `;
            table.appendChild(row);
        });

        container.appendChild(table);

        // Agregamos un bloque final mostrando cuanto dinero recaudo el museo en esas fechas
        const resumenDiv = document.createElement("div");
        resumenDiv.style.marginTop = "20px";
        resumenDiv.innerHTML = `
            <h4>Totales del Rango Seleccionado</h4>
            <p><strong>Monto Transaccionado:</strong> $${totalFacturas.toFixed(2)}</p>
            <p><strong>IVA Recaudado:</strong> $${totalIVA.toFixed(2)}</p>
            <p><strong>Ganancias del Museo:</strong> $${totalMuseo.toFixed(2)}</p>
        `;
        container.appendChild(resumenDiv);

    } catch (err) {
        console.error("Error conectando a Cassandra:", err);
        const container = document.getElementById("facturasContainer");
        container.innerHTML = "<p style='color:red;'>Error al conectar con la base de datos histórica. Verifica la consola.</p>";
    }
}