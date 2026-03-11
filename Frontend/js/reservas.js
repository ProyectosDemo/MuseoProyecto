document.addEventListener('DOMContentLoaded', cargarPagina);

const navMenu = document.getElementById("navMenu");
const loginButton = document.getElementById("loginButton");
const trabajadorId = localStorage.getItem("trabajadorId");
const trabajadorNombre = localStorage.getItem("trabajadorNombre");
const trabajadorAdmin = localStorage.getItem("trabajadorAdmin"); // "true" si es admin

async function cargarPagina() {
    header();
    await cargarReservas();
}

async function obtenerTelefonoCliente(id_cliente) {
    const query = `
        query FindCliente($id: ID!) {
            findCliente(id: $id) {
                telefono
            }
        }
    `;
    try {
        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query, variables: { id: id_cliente } })
        });
        const data = await res.json();
        const cliente = data?.data?.findCliente?.[0];
        return cliente?.telefono || "";
    } catch (err) {
        console.error("Error obteniendo teléfono del cliente:", err);
        return "";
    }
}

async function cargarReservas() {
    const query = `
        query {
            Ordenes(limit: 100, offset: 0) {
                id
                id_obra
                id_cliente
                id_trabajador
                fecha
                status
                direccion
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
        const ordenes = data?.data?.Ordenes || [];
        const tbody = document.querySelector("#tablaReservas tbody");
        tbody.innerHTML = "";

        for (let index = 0; index < ordenes.length; index++) {
            const orden = ordenes[index];
            const telefono = await obtenerTelefonoCliente(orden.id_cliente);

            const tr = document.createElement("tr");
            if (orden.status === "CONCRETADA") tr.classList.add("concretada");

            let accionHtml = "";
            if (!orden.id_trabajador) {
                accionHtml = `<button data-index="${index}" class="asignarBtn">Asignar</button>`;
            } else if (orden.status === "CONCRETADA") {
                accionHtml = `<span>CONCRETADA</span>`;
            } else {
                accionHtml = `
                    <select data-index="${index}" class="statusSelect">
                        <option value="PENDIENTE" ${orden.status === "PENDIENTE" ? "selected" : ""}>PENDIENTE</option>
                        <option value="CONCRETADA" ${orden.status === "CONCRETADA" ? "selected" : ""}>CONCRETADA</option>
                        <option value="CANCELADA" ${orden.status === "CANCELADA" ? "selected" : ""}>CANCELADA</option>
                    </select>
                    <button data-index="${index}" class="updateStatusBtn">Actualizar</button>
                `;
            }

            let direccionHtml = "";
            if (orden.status === "CONCRETADA" && orden.direccion) {
                direccionHtml = `<div class="direccionText">Dirección: ${orden.direccion}</div>`;
            }

            tr.innerHTML = `
                <td>${orden.id_obra}</td>
                <td>${orden.id_cliente}</td>
                <td>${telefono}</td>
                <td>${orden.fecha.substring(0, 10)}</td>
                <td>${orden.status}</td>
                <td>${accionHtml}${direccionHtml}</td>
            `;
            tbody.appendChild(tr);
        }

        window.ordenesGlobal = ordenes;

        // Eventos de botones
        tbody.querySelectorAll(".asignarBtn").forEach(btn => {
            btn.addEventListener("click", () => asignarTrabajador(btn.dataset.index));
        });

        tbody.querySelectorAll(".updateStatusBtn").forEach(btn => {
            btn.addEventListener("click", async () => {
                const select = document.querySelector(`select[data-index='${btn.dataset.index}']`);
                const nuevoStatus = select.value;

                const orden = window.ordenesGlobal[btn.dataset.index];

                if (orden.status === "CONCRETADA") {
                    alert("No se puede cambiar el estado de una orden concretada");
                    select.value = "CONCRETADA";
                    return;
                }

                const tr = btn.closest("tr");
                tr.cells[4].textContent = nuevoStatus;

                // Si cambió a concretada, mostrar formulario dirección
                if (nuevoStatus === "CONCRETADA") {
                    mostrarFormularioDireccion(btn.dataset.index, orden);
                    return;
                }

                // Actualizamos en memoria y en servidor
                orden.status = nuevoStatus;
                await updateOrdenStatus(orden.id, nuevoStatus);
            });
        });

    } catch (err) {
        console.error("Error al cargar órdenes:", err);
    }

    if (!trabajadorId) {
        location.href = "login.html";
    }

    if (trabajadorAdmin === "true") {
        document.getElementById("menuLateral").hidden = false;
    }
}

function header() {
    loginButton.textContent = `Bienvenido, ${trabajadorNombre || "Trabajador"}`;
    loginButton.href = "#";

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
}

async function asignarTrabajador(index) {
    const orden = window.ordenesGlobal[index];
    const mutation = `
        mutation UpdateOrden($id: ID!, $id_trabajador: ID!) {
            updateOrden(input: { id: $id, id_trabajador: $id_trabajador }) {
                id
                id_trabajador
            }
        }
    `;
    try {
        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query: mutation, variables: { id: orden.id, id_trabajador: trabajadorId } })
        });
        const data = await res.json();
        if (!data.errors) cargarReservas();
    } catch (err) {
        console.error("Error al asignar trabajador:", err);
    }
}

function mostrarFormularioDireccion(index, orden) {
    const tr = document.querySelector(`#tablaReservas tbody tr:nth-child(${parseInt(index)+1})`);
    document.querySelectorAll(".direccion-row").forEach(row => row.remove());

    const formRow = document.createElement("tr");
    formRow.classList.add("direccion-row");
    formRow.innerHTML = `
        <td colspan="6">
            <form id="direccionForm${index}">
                <label>Dirección de entrega:</label>
                <input type="text" name="direccion" required>
                <button type="submit">Guardar Dirección</button>
                <button type="button" id="cancelDireccion${index}">Cancelar</button>
            </form>
        </td>
    `;
    tr.after(formRow);

    const form = document.getElementById(`direccionForm${index}`);
    form.addEventListener("submit", async (e) => {
        e.preventDefault();
        const direccion = form.direccion.value.trim();
        if (!direccion) return;

        await updateOrdenStatus(orden.id, "CONCRETADA", direccion);
        formRow.remove();

        tr.cells[4].textContent = "CONCRETADA";
        const tdAccion = tr.cells[5];
        tdAccion.innerHTML = `<div class="direccionText">Dirección: ${direccion}</div><span>CONCRETADA</span>`;
    });

    document.getElementById(`cancelDireccion${index}`).addEventListener("click", () => {
        formRow.remove();
        const select = tr.querySelector(".statusSelect");
        if (select) select.value = orden.status;
    });
}

async function updateOrdenStatus(id, status, direccion = null) {
    const mutation = `
        mutation UpdateOrden($id: ID!, $status: StatusOrden!, $direccion: String) {
            updateOrden(input: { id: $id, status: $status, direccion: $direccion }) {
                id
                status
                direccion
            }
        }
    `;
    try {
        await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query: mutation, variables: { id, status, direccion } })
        });
    } catch (err) {
        console.error("Error al actualizar orden:", err);
    }
}