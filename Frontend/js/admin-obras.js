document.addEventListener('DOMContentLoaded', cargar);

async function cargar() {
    const navMenu = document.getElementById("navMenu");
const loginButton = document.getElementById("loginButton");
const trabajadorId = localStorage.getItem("trabajadorId");
const trabajadorNombre = localStorage.getItem("trabajadorNombre");
const trabajadorAdmin = localStorage.getItem("trabajadorAdmin"); 

if (!trabajadorId) {
    console.log("Debe iniciar sesión como trabajador para usar esta página");
    location.href = "login.html";
}

if (trabajadorAdmin === "true") {
    document.getElementById("menuLateral").hidden = false;
}

loginButton.textContent = `Bienvenido, ${trabajadorNombre || "Trabajador"}`;
loginButton.href = "#";

const logoutBtn = document.createElement("a");
logoutBtn.href = "#";
logoutBtn.textContent = "Cerrar Sesión";
logoutBtn.addEventListener("click", () => {
    localStorage.removeItem("trabajadorNombre");
    localStorage.removeItem("trabajadorId");
    localStorage.removeItem("trabajadorAdmin");
    location.reload();
});
navMenu.appendChild(logoutBtn);

const contenido = document.getElementById("contenidoClientes");
const endpoint = "http://localhost:8080/query";

function limpiar() {
    contenido.innerHTML = "";
}

async function fetchGraphQL(query, variables = {}) {
    const res = await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query, variables })
    });
    return await res.json();
}

// CREAR
document.getElementById("btnCrear").addEventListener("click", () => {
    limpiar();

    contenido.innerHTML = `
        <form id="formCrear" class="form-container">

            <input type="text" id="nombre" placeholder="Nombre" required>
            <input type="text" id="idArtista" placeholder="ID Artista" required>
            <input type="text" id="idGenero" placeholder="ID Género" required>
            <input type="number" id="precio" placeholder="Precio" required>

            <input type="date" id="fechaCreacion" required>

            <select id="status" required>
                <option value="DISPONIBLE">DISPONIBLE</option>
                <option value="RESERVADA">RESERVADA</option>
                <option value="VENDIDA">VENDIDA</option>
            </select>

            <input type="text" id="foto" placeholder="Foto URL" required>

            <label style="margin-top:10px;">
                <input type="checkbox" id="esEscultura">
                Es escultura
            </label>

            <div id="camposEscultura" style="display:none;">
                <input type="text" id="material" placeholder="Material">
                <input type="number" id="peso" placeholder="Peso">
                <input type="text" id="dimensiones" placeholder="Dimensiones">
            </div>

            <div class="form-buttons">
                <button type="button" id="cancelarCrear" class="btn-cancel">Cancelar</button>
                <button type="submit" class="btn-primary">Guardar</button>
            </div>

        </form>
    `;

    const checkbox = document.getElementById("esEscultura");
    const camposEscultura = document.getElementById("camposEscultura");

    // mostrar u ocultar campos
    checkbox.addEventListener("change", () => {
        camposEscultura.style.display = checkbox.checked ? "block" : "none";
    });

    document.getElementById("formCrear").addEventListener("submit", async e => {
        e.preventDefault();

        // 1️⃣ crear obra
        const mutationObra = `
        mutation ($input: NewObra!) {
            createObra(input: $input) {
                id
                nombre
            }
        }`;

        const variablesObra = {
            input: {
                nombre: document.getElementById("nombre").value,
                id_artista: document.getElementById("idArtista").value,
                id_genero: document.getElementById("idGenero").value,
                precio: parseInt(document.getElementById("precio").value),
                fecha_creacion: document.getElementById("fechaCreacion").value,
                status: document.getElementById("status").value,
                foto: document.getElementById("foto").value
            }
        };

        const resultObra = await fetchGraphQL(mutationObra, variablesObra);

        if (resultObra.errors) {
            alert("Error al crear obra");
            return;
        }

        const obraID = resultObra.data.createObra.id;

        // 2️⃣ si es escultura → crear escultura
        if (checkbox.checked) {

            const mutationEscultura = `
            mutation ($input: NewEscultura!) {
                createEscultura(input: $input) {
                    id_obra
                }
            }`;

            const variablesEscultura = {
                input: {
                    id_obra: obraID,
                    material: document.getElementById("material").value,
                    peso: parseInt(document.getElementById("peso").value),
                    dimensiones: document.getElementById("dimensiones").value
                }
            };

            const resultEscultura = await fetchGraphQL(mutationEscultura, variablesEscultura);

            if (resultEscultura.errors) {
                alert("La obra se creó pero hubo error con la escultura");
                return;
            }
        }

        alert("Obra creada correctamente");
        limpiar();
    });

    document.getElementById("cancelarCrear").addEventListener("click", limpiar);
});

// OBTENER
document.getElementById("btnObtener").addEventListener("click", async () => {
    limpiar();

    const query = `
        query ObtenerObras($limit: Int, $offset: Int) {
            Obras(limit: $limit, offset: $offset) {
                id
                nombre
                id_artista
                id_genero
                precio
                fecha_creacion
                status
            }
        }
    `;

    const data = await fetchGraphQL(query, { limit: 100, offset: 0 });
    const obras = data?.data?.Obras;

    if (!obras || obras.length === 0) {
        contenido.innerHTML = "<p>No hay obras registradas.</p>";
        return;
    }

    // Generar HTML con tarjetas neutras en rejilla de 3 columnas
    let html = '<div class="grid-container"><h3>Obras Registradas</h3>';

    obras.forEach(o => {
        html += `<div class="grid-item">
                    <h4>${o.nombre}</h4>
                    <p><strong>Artista ID:</strong> ${o.id_artista}</p>
                    <p><strong>Género ID:</strong> ${o.id_genero}</p>
                    <p><strong>Precio:</strong> ${o.precio}</p>
                    <p><strong>Fecha de Creación:</strong> ${o.fecha_creacion}</p>
                    <p><strong>Status:</strong> ${o.status}</p>`;
        html += `</div>`;
    });

    html += '</div>';
    contenido.innerHTML = html;
});
// ELIMINAR
document.getElementById("btnEliminar").addEventListener("click", async () => {
    limpiar();
    const query = `query ($limit: Int, $offset: Int) { Obras(limit: $limit, offset: $offset) { id nombre } }`;
    const data = await fetchGraphQL(query, { limit: 100, offset: 0 });
    const obras = data?.data?.Obras;
    if (!obras || obras.length === 0) { contenido.innerHTML = "<p>No hay obras para eliminar.</p>"; return; }

    const opciones = obras.map(o => `<option value="${o.id}">${o.nombre}</option>`).join("");
    contenido.innerHTML = `
        <select id="selectObra" class="select-box">${opciones}</select>
        <div class="form-buttons">
            <button id="cancelarEliminar" class="btn-cancel">Cancelar</button>
            <button id="confirmarEliminar" class="btn-delete">Eliminar</button>
        </div>
    `;
    document.getElementById("confirmarEliminar").addEventListener("click", async () => {
        const id = document.getElementById("selectObra").value;
        if (!id) return alert("Selecciona una obra para eliminar");
        const mutation = `mutation KillObra($id: ID!) { killObra(id: $id) }`;
        const result = await fetchGraphQL(mutation, { id });
        if (result.errors) alert("Error al eliminar obra");
        else if (result.data.killObra) { alert("Obra eliminada correctamente"); limpiar(); }
        else alert("No se pudo eliminar la obra");
    });
    document.getElementById("cancelarEliminar").addEventListener("click", limpiar);
});

// EDITAR
// EDITAR
document.getElementById("btnEditar").addEventListener("click", async () => {
    limpiar();

    const query = `
    query ($limit: Int, $offset: Int) {
        Obras(limit: $limit, offset: $offset) {
            id
            nombre
            id_artista
            id_genero
            precio
            fecha_creacion
            status
            foto
        }
    }`;

    const data = await fetchGraphQL(query, { limit: 100, offset: 0 });
    const obras = data?.data?.Obras;

    if (!obras || obras.length === 0) {
        contenido.innerHTML = "<p>No hay obras para editar.</p>";
        return;
    }

    const opciones = obras.map(o => `<option value="${o.id}">${o.nombre}</option>`).join("");

    contenido.innerHTML = `
        <select id="selectObra" class="select-box">${opciones}</select>
        <form id="updateForm" class="form-container"></form>
    `;

    const form = document.getElementById("updateForm");

    async function obtenerEscultura(idObra) {

        const query = `
        query ($id: ID!) {
            findEscultura(id_obra: $id) {
                id_obra
                material
                peso
                dimensiones
            }
        }`;

        const data = await fetchGraphQL(query, { id: idObra });

        if (data?.data?.findEscultura)
            return data.data.findEscultura;

        return null;
    }

    async function renderForm(obra) {

        if (!obra) return;

        const escultura = await obtenerEscultura(obra.id);

        let camposEscultura = "";

        if (escultura) {
            camposEscultura = ` 
                <input type="text" id="material" value="${escultura.material || ''}" placeholder="Material">
                <input type="number" id="peso" value="${escultura.peso || 0}" placeholder="Peso">
                <input type="text" id="dimensiones" value="${escultura.dimensiones || ''}" placeholder="Dimensiones">
            `;
        }

        form.innerHTML = `
            <input type="text" id="nombre" value="${obra.nombre || ''}" required>
            <input type="number" id="precio" value="${obra.precio || 0}" required>
            <input type="text" id="fechaCreacion" value="${obra.fecha_creacion || ''}" required>

            ${camposEscultura}

            <div class="form-buttons">
                <button type="button" id="cancelarEditar" class="btn-cancel">Cancelar</button>
                <button type="submit" class="btn-primary">Actualizar</button>
            </div>
        `;

        form.onsubmit = async e => {
            e.preventDefault();

            // actualizar obra
            const mutationObra = `
            mutation updateObra($input: UpdateObra!) {
                updateObra(input: $input) { id nombre }
            }`;

            const variablesObra = {
                input: {
                    id: obra.id,
                    nombre: document.getElementById("nombre").value,
                    precio: parseInt(document.getElementById("precio").value),
                    fecha_creacion: document.getElementById("fechaCreacion").value
                }
            };

            const resultObra = await fetchGraphQL(mutationObra, variablesObra);

            if (resultObra.errors) {
                alert("Error al actualizar obra");
                return;
            }

            // actualizar escultura si existe
            if (escultura) {

                const mutationEscultura = `
                mutation updateEscultura($input: UpdateEscultura!) {
                    updateEscultura(input: $input) { id_obra }
                }`;

                const variablesEscultura = {
                    input: {
                        id_obra: obra.id,
                        material: document.getElementById("material").value,
                        peso: parseInt(document.getElementById("peso").value),
                        dimensiones: document.getElementById("dimensiones").value
                    }
                };

                const resultEscultura = await fetchGraphQL(mutationEscultura, variablesEscultura);

                if (resultEscultura.errors) {
                    alert("La obra se actualizó pero hubo error con la escultura");
                    return;
                }
            }

            alert("Obra actualizada correctamente");
            limpiar();
        };

        document.getElementById("cancelarEditar").addEventListener("click", limpiar);
    }

    renderForm(obras[0]);

    document.getElementById("selectObra").addEventListener("change", e => {
        const obra = obras.find(o => o.id === e.target.value);
        renderForm(obra);
    });
});
}