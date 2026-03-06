// logica de consultas

document.addEventListener('DOMContentLoaded', cargarOrdenesConcretadas);

const navMenu = document.getElementById("navMenu");
const loginButton = document.getElementById("loginButton");
const trabajadorId = localStorage.getItem("trabajadorId");
const trabajadorNombre = localStorage.getItem("trabajadorNombre");
const logoutBtn = document.createElement("a");

async function cargarConsultas() {
    if (!trabajadorId) location.href = "login.html";
}

function mostrarPanel(panelId) {
    const panels = document.querySelectorAll(".main-content .panel");
    panels.forEach(p => p.hidden = true);
    document.getElementById(panelId).hidden = false;
}

async function cargarOrdenesConcretadas() {
    try {
        const query = `query { Ordenes(limit: 100, offset: 0) { id id_obra id_cliente id_trabajador fecha status } }`;
        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query })
        });
        const data = await res.json();
        const ordenes = data?.data?.Ordenes.filter(o => o.status === "CONCRETADA") || [];
        const container = document.getElementById("ordenesContainer");
        container.innerHTML = "";

        for (const orden of ordenes) {
            const obraQuery = `query { findObra(id: "${orden.id_obra}") { id nombre precio } }`;
            const obraRes = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: obraQuery })
            });
            const obraData = await obraRes.json();
            const obra = obraData?.data?.findObra[0];

            const clienteQuery = `query { findCliente(id: "${orden.id_cliente}") { id nombre } }`;
            const clienteRes = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: clienteQuery })
            });
            const clienteData = await clienteRes.json();
            const cliente = clienteData?.data?.findCliente[0];

            const wrapper = document.createElement("div");
            wrapper.className = "orden-factura-wrapper";

            const ordenCard = document.createElement("div");
            ordenCard.className = "orden-card";
            ordenCard.innerHTML = `
                <p><strong>ID Orden:</strong> ${orden.id}</p>
                <p><strong>Obra:</strong> ${obra?.nombre || "N/D"}</p>
                <p><strong>Cliente:</strong> ${cliente?.nombre || "N/D"}</p>
                <p><strong>Fecha:</strong> ${orden.fecha}</p>
                <p><strong>Precio:</strong> $${obra?.precio || 0}</p>
                <button class="facturaBtn">Generar Factura</button>
            `;

            const facturaCard = document.createElement("div");
            facturaCard.className = "factura-card";
            facturaCard.innerHTML = `<p>Factura aparecerá aquí al seleccionar la orden.</p>`;

            wrapper.appendChild(ordenCard);
            wrapper.appendChild(facturaCard);
            container.appendChild(wrapper);

            ordenCard.querySelector(".facturaBtn").addEventListener("click", () => {
                generarFacturaEnCard(orden, obra, cliente, facturaCard);
            });
        }
    } catch (err) {
        console.error("Error al cargar ordenes concretadas:", err);
    }
}

function generarFacturaEnCard(orden, obra, cliente, card) {
    const IVA = 0.21;
    const porcentajeMuseo = 0.10;
    const precio = obra.precio;
    const iva = precio * IVA;
    const gananciaMuseo = precio * porcentajeMuseo;
    const total = precio + iva;

    card.innerHTML = `
        <p><strong>Factura ID:</strong> ${orden.id}</p>
        <p><strong>Cliente:</strong> ${cliente.nombre}</p>
        <p><strong>Obra:</strong> ${obra.nombre}</p>
        <p><strong>Precio:</strong> $${precio.toFixed(2)}</p>
        <p><strong>IVA (21%):</strong> $${iva.toFixed(2)}</p>
        <p><strong>Total:</strong> $${total.toFixed(2)}</p>
        <p><strong>Ganancia Museo:</strong> $${gananciaMuseo.toFixed(2)}</p>
        <p><strong>Fecha:</strong> ${orden.fecha}</p>
    `;
}

async function cargarResumenMembresias(fechaInicio, fechaFin) {
    try {
        const query = `query {
            Membresias(limit: 1000, offset: 0) {
                id_membresia
                fecha
                cliente { nombre }
                tarjeta { numero_tarjeta }
            }
        }`;
        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query })
        });
        const data = await res.json();
        const membresias = data?.data?.Membresias || [];
        const filtradas = membresias.filter(m => {
            const fechaM = new Date(m.fecha);
            const inicio = new Date(fechaInicio);
            const fin = new Date(fechaFin);
            return fechaM >= inicio && fechaM <= fin;
        });

        const container = document.getElementById("resumenMembresias");
        container.innerHTML = "<h3>Resumen de Membresías</h3>";
        if (filtradas.length === 0) {
            container.innerHTML += "<p>No hay membresías en este periodo.</p>";
            return;
        }

        const table = document.createElement("table");
        table.innerHTML = "<tr><th>Cliente</th><th>Fecha de Membresía</th></tr>";
        filtradas.forEach(m => {
            const row = document.createElement("tr");
            row.innerHTML = `<td>${m.cliente.nombre}</td><td>${new Date(m.fecha).toLocaleDateString()}</td>`;
            table.appendChild(row);
        });
        container.appendChild(table);
    } catch (err) {
        console.error("Error al cargar resumen de membresías:", err);
    }
}

async function cargarFacturasPorPeriodo(fechaInicio, fechaFin) {
    try {
        const query = `query { Ordenes(limit: 1000, offset: 0) { id id_obra id_cliente fecha status } }`;
        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query })
        });
        const data = await res.json();
        const ordenes = data?.data?.Ordenes.filter(o => o.status === "CONCRETADA") || [];

        const filtradas = ordenes.filter(o => {
            const fechaO = new Date(o.fecha);
            const inicio = new Date(fechaInicio);
            const fin = new Date(fechaFin);
            return fechaO >= inicio && fechaO <= fin;
        });

        const container = document.getElementById("reportesContainer");
        container.innerHTML = "<h3>Facturas del periodo</h3>";
        if (filtradas.length === 0) {
            container.innerHTML += "<p>No hay facturas en este periodo.</p>";
            return;
        }

        const table = document.createElement("table");
        table.innerHTML = "<tr><th>ID Orden</th><th>Cliente</th><th>Obra</th><th>Precio</th><th>IVA</th><th>Total</th><th>Ganancia Museo</th><th>Fecha</th></tr>";

        for (const orden of filtradas) {
            const obraRes = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: `query { findObra(id: "${orden.id_obra}") { id nombre precio } }` })
            });
            const obraData = await obraRes.json();
            const obra = obraData?.data?.findObra[0];

            const clienteRes = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: `query { findCliente(id: "${orden.id_cliente}") { id nombre } }` })
            });
            const clienteData = await clienteRes.json();
            const cliente = clienteData?.data?.findCliente[0];

            const precio = obra?.precio || 0;
            const IVA = 0.21;
            const porcentajeMuseo = 0.10;
            const iva = precio * IVA;
            const gananciaMuseo = precio * porcentajeMuseo;
            const total = precio + iva;

            const row = document.createElement("tr");
            row.innerHTML = `
                <td>${orden.id}</td>
                <td>${cliente?.nombre || "N/D"}</td>
                <td>${obra?.nombre || "N/D"}</td>
                <td>$${precio.toFixed(2)}</td>
                <td>$${iva.toFixed(2)}</td>
                <td>$${total.toFixed(2)}</td>
                <td>$${gananciaMuseo.toFixed(2)}</td>
                <td>${new Date(orden.fecha).toLocaleDateString()}</td>
            `;
            table.appendChild(row);
        }
        container.appendChild(table);
    } catch (err) {
        console.error("Error al cargar facturas por periodo:", err);
    }
}
