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

document.getElementById("btnFiltrarMembresias").addEventListener("click", async () => {
    const desde = document.getElementById("fechaInicio").value;
    const hasta = document.getElementById("fechaFin").value;
    if (!desde || !hasta) return alert("Selecciona ambas fechas.");
    cargarResumenMembresias(desde, hasta);
});

async function cargarResumenMembresias(fechaInicio, fechaFin) {
    const query = `query { Membresias(limit: 1000, offset: 0) { id_membresia fecha cliente { nombre } tarjeta { numero_tarjeta } } }`;
    const res = await fetch("http://localhost:8080/query", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query })
    });
    const data = await res.json();
    const membresias = data?.data?.Membresias || [];
    const filtradas = membresias.filter(m => {
        const fechaM = new Date(m.fecha);
        return fechaM >= new Date(fechaInicio) && fechaM <= new Date(fechaFin);
    });

    const container = document.getElementById("resumenMembresias");
    container.innerHTML = "<h3>Resumen de Membresías</h3>";
    if (!filtradas.length) return container.innerHTML += "<p>No hay membresías en este periodo.</p>";

    const table = document.createElement("table");
    table.innerHTML = "<tr><th>Cliente</th><th>Fecha de Membresía</th></tr>";
    filtradas.forEach(m => {
        const row = document.createElement("tr");
        row.innerHTML = `<td>${m.cliente.nombre}</td><td>${new Date(m.fecha).toLocaleDateString()}</td>`;
        table.appendChild(row);
    });
    container.appendChild(table);
}

document.getElementById("btnFiltrarObras").addEventListener("click", async () => {
    const desde = document.getElementById("fechaObrasInicio").value;
    const hasta = document.getElementById("fechaObrasFin").value;
    if (!desde || !hasta) return alert("Selecciona ambas fechas.");
    cargarObrasVendidas(desde, hasta);
});

async function cargarObrasVendidas(fechaInicio, fechaFin) {
    const query = `query { Ordenes(limit: 1000, offset: 0) { id id_obra fecha status } }`;
    const res = await fetch("http://localhost:8080/query", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query })
    });
    const data = await res.json();
    const ordenes = data?.data?.Ordenes.filter(o => o.status === "CONCRETADA") || [];
    const obrasVendidasMap = {};

    for (const orden of ordenes) {
        const fechaOrden = new Date(orden.fecha);
        if (fechaOrden < new Date(fechaInicio) || fechaOrden > new Date(fechaFin)) continue;

        const obraQuery = `query { findObra(id: "${orden.id_obra}") { nombre precio } }`;
        const obraRes = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query: obraQuery })
        });
        const obraData = await obraRes.json();
        const obra = obraData?.data?.findObra[0];
        if (!obra) continue;

        if (!obrasVendidasMap[obra.nombre]) obrasVendidasMap[obra.nombre] = { cantidad: 0, total: 0, precio: obra.precio };
        obrasVendidasMap[obra.nombre].cantidad++;
        obrasVendidasMap[obra.nombre].total += obra.precio;
    }

    const container = document.getElementById("obrasContainer");
    container.innerHTML = "<h3>Obras Vendidas por Fecha</h3>";
    const table = document.createElement("table");
    table.innerHTML = "<tr><th>Obra</th><th>Cantidad Vendida</th><th>Precio Unitario</th><th>Total Vendido</th></tr>";

    for (const [nombre, info] of Object.entries(obrasVendidasMap)) {
        const row = document.createElement("tr");
        row.innerHTML = `<td>${nombre}</td><td>${info.cantidad}</td><td>$${info.precio.toFixed(2)}</td><td>$${info.total.toFixed(2)}</td>`;
        table.appendChild(row);
    }
    container.appendChild(table);
}

document.getElementById("btnFiltrarFacturas").addEventListener("click", async () => {
    const desde = document.getElementById("fechaFacturasInicio").value;
    const hasta = document.getElementById("fechaFacturasFin").value;
    if (!desde || !hasta) return alert("Selecciona ambas fechas.");
    cargarResumenFacturas(desde, hasta);
});

async function cargarResumenFacturas(fechaInicio, fechaFin) {
    const query = `query { Ordenes(limit: 1000, offset: 0) { id id_obra id_cliente fecha status } }`;
    const res = await fetch("http://localhost:8080/query", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query })
    });
    const data = await res.json();
    const ordenes = data?.data?.Ordenes.filter(o => o.status === "CONCRETADA") || [];
    const filtradas = [];

    for (const orden of ordenes) {
        const fechaOrden = new Date(orden.fecha);
        if (fechaOrden < new Date(fechaInicio) || fechaOrden > new Date(fechaFin)) continue;

        const obraQuery = `query { findObra(id: "${orden.id_obra}") { nombre precio } }`;
        const obraRes = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query: obraQuery })
        });
        const obraData = await obraRes.json();
        const obra = obraData?.data?.findObra[0];
        if (!obra) continue;

        const clienteQuery = `query { findCliente(id: "${orden.id_cliente}") { nombre } }`;
        const clienteRes = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query: clienteQuery })
        });
        const clienteData = await clienteRes.json();
        const cliente = clienteData?.data?.findCliente[0];

        filtradas.push({ orden, obra, cliente });
    }

    const container = document.getElementById("facturasContainer");
    container.innerHTML = "<h3>Resumen de Facturas</h3>";
    if (!filtradas.length) return container.innerHTML += "<p>No hay facturas en este periodo.</p>";

    const IVA = 0.10;
    const porcentajeMuseo = 0.10;
    let totalFacturas = 0, totalIVA = 0, totalMuseo = 0;

    const table = document.createElement("table");
    table.innerHTML = "<tr><th>ID Orden</th><th>Obra</th><th>Cliente</th><th>Precio</th><th>IVA</th><th>Total</th><th>Ganancia Museo</th><th>Fecha</th></tr>";

    filtradas.forEach(f => {
        const precio = f.obra.precio;
        const iva = precio * IVA;
        const gananciaMuseo = precio * porcentajeMuseo;
        const total = precio + iva;
        totalFacturas += total;
        totalIVA += iva;
        totalMuseo += gananciaMuseo;

        const row = document.createElement("tr");
        row.innerHTML = `<td>${f.orden.id}</td><td>${f.obra.nombre}</td><td>${f.cliente?.nombre || "N/D"}</td>
                         <td>$${precio.toFixed(2)}</td><td>$${iva.toFixed(2)}</td><td>$${total.toFixed(2)}</td>
                         <td>$${gananciaMuseo.toFixed(2)}</td><td>${new Date(f.orden.fecha).toLocaleDateString()}</td>`;
        table.appendChild(row);
    });

    container.appendChild(table);

}