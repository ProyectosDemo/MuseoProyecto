//logica de la pagina administracion de clientes

// Llamar a la función para cargar el detalle del artista al cargar la página
document.addEventListener('DOMContentLoaded', cargar);

const navMenu = document.getElementById("navMenu");
const loginButton = document.getElementById("loginButton");
const trabajadorId = localStorage.getItem("trabajadorId");
const trabajadorNombre = localStorage.getItem("trabajadorNombre");
const trabajadorAdmin = localStorage.getItem("trabajadorAdmin"); // "true" si es admin
const logoutBtn = document.createElement("a");
const contenido = document.getElementById("contenidoClientes");
const endpoint = "http://localhost:8080/query";

async function cargar() {

    if (!trabajadorId) {
        console.log("Debe iniciar sesión como trabajador para usar esta página");
        location.href = "login.html";
    }

    if (trabajadorAdmin === "true") {
        document.getElementById("menuLateral").hidden = false;
    }

    loginButton.textContent = `Bienvenido, ${trabajadorNombre || "Trabajador"}`;
    loginButton.href = "#";

    logoutBtn.href = "#";
    logoutBtn.textContent = "Cerrar Sesión";
    logoutBtn.addEventListener("click", () => {
        localStorage.removeItem("trabajadorNombre");
        localStorage.removeItem("trabajadorId");
        localStorage.removeItem("trabajadorAdmin");
        location.reload();
    });
    navMenu.appendChild(logoutBtn);

    document.getElementById("btnCrear").addEventListener("click", () => {
        limpiar();

        contenido.innerHTML = `
        <form id="formCrear" class="form-container">
            <input type="text" id="nombre" placeholder="Nombre" required>
            <input type="email" id="email" placeholder="Email" required>
            <input type="text" id="telefono" placeholder="Teléfono" required>
            <input type="text" id="login" placeholder="Login" required>
            <input type="text" id="password" placeholder="Contraseña" required>
            <input type="text" id="codigoSeguridad" placeholder="Código Seguridad" required>
            <div class="form-buttons">
                <button type="button" id="cancelarCrear" class="btn-cancel">Cancelar</button>
                <button type="submit" class="btn-primary">Guardar</button>
            </div>
        </form>
    `;

        document.getElementById("formCrear").addEventListener("submit", async e => {
            e.preventDefault();

            const query = `
            mutation ($input: NewCliente!) {
                createCliente(input: $input) {
                    id
                    nombre
                }
            }
        `;

            const variables = {
                input: {
                    nombre: document.getElementById("nombre").value,
                    email: document.getElementById("email").value,
                    telefono: document.getElementById("telefono").value,
                    login: document.getElementById("login").value,
                    password: document.getElementById("password").value,
                    codigo_seguridad: document.getElementById("codigoSeguridad").value
                }
            };

            const data = await fetchGraphQL(query, variables);
            console.log(data);

            if (data.errors) {
                console.error(data.errors);
                alert("Error al crear cliente");
            } else {
                alert("Cliente creado correctamente");
                limpiar();
            }
        });

        document.getElementById("cancelarCrear").addEventListener("click", limpiar);
    });
    document.getElementById("btnObtener").addEventListener("click", async () => {
        limpiar();

        const query = `
        query ObtenerClientes($limit: Int, $offset: Int) {
            Clientes(limit: $limit, offset: $offset) {
                id
                nombre
                email
                telefono
                login
                codigo_seguridad
            }
        }
    `;

        const variables = { limit: 100, offset: 0 };

        const data = await fetchGraphQL(query, variables);
        console.log(data);

        const clientes = data?.data?.Clientes;

        if (!clientes || clientes.length === 0) {
            contenido.innerHTML = "<p>No hay clientes registrados.</p>";
            return;
        }

        // Generar HTML con tarjetas neutras en rejilla
        let html = '<div class="grid-container"><h3>Clientes Registrados</h3>';

        clientes.forEach(c => {
            html += `<div class="grid-item">
                    <h4>${c.nombre}</h4>
                    <p><strong>Email:</strong> ${c.email}</p>
                    <p><strong>Teléfono:</strong> ${c.telefono}</p>
                    <p><strong>Login:</strong> ${c.login}</p>
                    <p><strong>Código de Seguridad:</strong> ${c.codigo_seguridad}</p>
                 </div>`;
        });

        html += '</div>';

        contenido.innerHTML = html;
    });

    document.getElementById("btnEliminar").addEventListener("click", async () => {
        limpiar();

        const query = `
        query ($limit: Int, $offset: Int) {
            Clientes(limit: $limit, offset: $offset) {
                id
                nombre
            }
        }
    `;
        const variables = { limit: 100, offset: 0 };
        const data = await fetchGraphQL(query, variables);
        const clientes = data?.data?.Clientes;

        if (!clientes || clientes.length === 0) {
            contenido.innerHTML = "<p>No hay clientes para eliminar.</p>";
            return;
        }

        const opciones = clientes.map(c => `<option value="${c.id}">${c.nombre}</option>`).join("");

        contenido.innerHTML = `
        <select id="selectCliente" class="select-box">${opciones}</select>
        <div class="form-buttons">
            <button id="cancelarEliminar" class="btn-cancel">Cancelar</button>
            <button id="confirmarEliminar" class="btn-delete">Eliminar</button>
        </div>
    `;

        document.getElementById("confirmarEliminar").addEventListener("click", async () => {
            const id = document.getElementById("selectCliente").value;
            if (!id) return alert("Selecciona un cliente para eliminar");

            const mutation = `
            mutation KillCliente($id: ID!) {
                killCliente(id: $id)
            }
        `;
            const variables = { id };

            const result = await fetchGraphQL(mutation, variables);
            if (result.errors) {
                console.error(result.errors);
                alert("Error al eliminar cliente");
            } else if (result.data.killCliente) {
                alert("Cliente eliminado correctamente");
                limpiar();
            } else {
                alert("No se pudo eliminar el cliente");
            }
        });

        document.getElementById("cancelarEliminar").addEventListener("click", limpiar);
    });

    document.getElementById("btnEditar").addEventListener("click", async () => {
        limpiar();

        const query = `
        query ($limit: Int, $offset: Int) {
            Clientes(limit: $limit, offset: $offset) {
                id
                nombre
                email
                telefono
                login
                password
                codigo_seguridad
            }
        }
    `;
        const variables = { limit: 100, offset: 0 };
        const data = await fetchGraphQL(query, variables);
        const clientes = data?.data?.Clientes;

        if (!clientes || clientes.length === 0) {
            contenido.innerHTML = "<p>No hay clientes para editar.</p>";
            return;
        }

        const opciones = clientes.map(c => `<option value="${c.id}">${c.nombre}</option>`).join("");

        contenido.innerHTML = `
        <select id="selectCliente" class="select-box">${opciones}</select>
        <form id="updateForm" class="form-container"></form>
    `;

        const form = document.getElementById("updateForm");
        form.classList.add("edit-form");


        renderForm(clientes[0]);

        document.getElementById("selectCliente").addEventListener("change", e => {
            const cliente = clientes.find(c => c.id === e.target.value);
            renderForm(cliente);
        });
    });
}

function renderForm(cliente) {
    if (!cliente) return;
    const form = document.getElementById("updateForm");
    if (!form) return;
    form.innerHTML = `
            <input type="text" id="nombre" value="${cliente.nombre || ''}" required>
            <input type="email" id="email" value="${cliente.email || ''}" required>
            <input type="text" id="telefono" value="${cliente.telefono || ''}" required>
            <input type="text" id="login" value="${cliente.login || ''}" required>
            <input type="text" id="password" value="${cliente.password || ''}" required>
            <input type="text" id="codigo_seguridad" value="${cliente.codigo_seguridad || ''}" required>
            <div class="form-buttons">
                <button type="button" id="cancelarEditar" class="btn-cancel">Cancelar</button>
                <button type="submit" class="btn-primary">Actualizar</button>
            </div>
        `;

    form.addEventListener("submit", async e => {
        e.preventDefault();

        const mutation = `
                mutation updateCliente($input: UpdateCliente!) {
                    updateCliente(input: $input) {
                        id
                        nombre
                        email
                        telefono
                        login
                        password
                        codigo_seguridad
                    }
                }
            `;

        const variables = {
            input: {
                id: cliente.id,
                nombre: document.getElementById("nombre").value,
                email: document.getElementById("email").value,
                telefono: document.getElementById("telefono").value,
                login: document.getElementById("login").value,
                password: document.getElementById("password").value,
                codigo_seguridad: document.getElementById("codigo_seguridad").value
            }
        };

        const result = await fetchGraphQL(mutation, variables);
        if (result.errors) {
            console.error(result.errors);
            alert("Error al actualizar cliente");
        } else {
            alert("Cliente actualizado correctamente");
            limpiar();
        }
    });

    document.getElementById("cancelarEditar").addEventListener("click", limpiar);
}

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
