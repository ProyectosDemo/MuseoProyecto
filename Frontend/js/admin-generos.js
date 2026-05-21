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

const contenido = document.getElementById("contenidoGeneros");
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

document.getElementById("btnCrear").addEventListener("click", () => {
    limpiar();
    contenido.innerHTML = `
        <form id="formCrear" class="form-container">
            <input type="text" id="nombre" placeholder="Nombre" required>
            <div class="form-buttons">
                <button type="button" id="cancelarCrear" class="btn-cancel">Cancelar</button>
                <button type="submit" class="btn-primary">Guardar</button>
            </div>
        </form>
    `;
    document.getElementById("formCrear").addEventListener("submit", async e => {
        e.preventDefault();
        
        // CORRECCION: Declaramos la variable de la mutacion y la pasamos correctamente
        const query = `
            mutation CrearNuevoGenero($nombre: String!) {
                createGenero(nombre: $nombre) {
                    id
                    nombre
                }
            }
        `;
        
        const variables = { nombre: document.getElementById("nombre").value };
        const data = await fetchGraphQL(query, variables);
        console.log(data);
        
        if (data.errors) {
            console.error(data.errors);
            alert("Error al crear genero: " + data.errors[0].message);
        } else {
            alert("Genero creado correctamente");
            limpiar();
        }
    });
    document.getElementById("cancelarCrear").addEventListener("click", limpiar);
});


document.getElementById("btnObtener").addEventListener("click", async () => {
    limpiar();

    const query = `
        query ObtenerGeneros($limit: Int, $offset: Int) {
            Genero(limit: $limit, offset: $offset) {
                id
                nombre
            }
        }
    `;
    const variables = { limit: 100, offset: 0 };
    const data = await fetchGraphQL(query, variables);
    console.log(data);

    const generos = data?.data?.Genero;
    if (!generos || generos.length === 0) {
        contenido.innerHTML = "<p>No hay géneros registrados.</p>";
        return;
    }

    // Generar HTML con tarjetas neutras en rejilla
    let html = '<div class="grid-container"><h3>Géneros Registrados</h3>';

    generos.forEach(g => {
        html += `<div class="grid-item">
                    <h4>${g.nombre}</h4>
                 </div>`;
    });

    html += '</div>';

    contenido.innerHTML = html;
});

document.getElementById("btnEditar").addEventListener("click", async () => {
    limpiar();
    const query = `
        query ($limit: Int, $offset: Int) {
            Genero(limit: $limit, offset: $offset) {
                id
                nombre
            }
        }
    `;
    const variables = { limit: 100, offset: 0 };
    const data = await fetchGraphQL(query, variables);
    const generos = data?.data?.Genero;
    if (!generos || generos.length === 0) {
        contenido.innerHTML = "<p>No hay géneros para editar.</p>";
        return;
    }
    const opciones = generos.map(g => `<option value="${g.id_genero}">${g.nombre}</option>`).join("");
    contenido.innerHTML = `
        <select id="selectGenero" class="select-box">${opciones}</select>
        <form id="updateForm" class="form-container"></form>
    `;
    const form = document.getElementById("updateForm");
    form.classList.add("edit-form");

    function renderForm(genero) {
        if (!genero) return;
        form.innerHTML = `
            <input type="text" id="nombre" value="${genero.nombre || ''}" required>
            <div class="form-buttons">
                <button type="button" id="cancelarEditar" class="btn-cancel">Cancelar</button>
                <button type="submit" class="btn-primary">Actualizar</button>
            </div>
        `;
        document.getElementById("updateForm").addEventListener("submit", async e => {
            e.preventDefault();
            const mutation = `
                mutation updateGenero($input: UpdateGenero!) {
                    updateGenero(input: $input) {
                        id_genero
                        nombre
                    }
                }
            `;
            const variables = {
                input: {
                    id: genero.id_genero,
                    nombre: document.getElementById("nombre").value
                }
            };
            const result = await fetchGraphQL(mutation, variables);
            if (result.errors) {
                console.error(result.errors);
                alert("Error al actualizar género");
            } else {
                alert("Género actualizado correctamente");
                limpiar();
            }
        });
        document.getElementById("cancelarEditar").addEventListener("click", limpiar);
    }

    renderForm(generos[0]);
    document.getElementById("selectGenero").addEventListener("change", e => {
        const genero = generos.find(g => g.id_genero === e.target.value);
        renderForm(genero);
    });
});
}