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

const contenido = document.getElementById("contenidoPreguntas");
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
// Mostrar ítems en tarjetas agrupadas
document.getElementById("btnObtener").addEventListener("click", async () => {
    limpiar();

    const query = `
        query ObtenerItems($limit: Int, $offset: Int) {
            Preguntas(limit: $limit, offset: $offset) {
                id
                pregunta
                respuesta
                cliente {
                    nombre
                }
            }
        }
    `;

    const variables = { limit: 100, offset: 0 };
    const data = await fetchGraphQL(query, variables);

    const items = data?.data?.Preguntas;
    if (!items || items.length === 0) {
        contenido.innerHTML = "<p>No hay registros.</p>";
        return;
    }

    // Agrupar por categoría/neutro (antes cliente)
    const groupedItems = {};
    items.forEach(i => {
        const group = i.cliente?.nombre || 'Desconocido';
        if (!groupedItems[group]) groupedItems[group] = [];
        groupedItems[group].push(i);
    });

    // Generar HTML neutral
    let html = '<div class="grid-container"><h3>Registros</h3>';
    for (const group in groupedItems) {
        html += `<div class="grid-item">
                    <h4>${group}</h4>`;
        groupedItems[group].forEach(i => {
            html += `<p><strong>Pregunta:</strong> ${i.pregunta}</p>
                     <p><strong>Respuesta:</strong> ${i.respuesta}</p>
                     <hr>`;
        });
        html += `</div>`;
    }
    html += '</div>';

    contenido.innerHTML = html;
});

}