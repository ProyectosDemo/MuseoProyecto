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

const contenido = document.getElementById("contenidoMembresias");
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
document.getElementById("btnObtener").addEventListener("click", async () => {
    limpiar();

    const query = `
        query ObtenerMembresias($limit: Int, $offset: Int) {
            Membresias(limit: $limit, offset: $offset) {
                id_membresia
                cliente {
                    nombre
                }
                id_tarjeta
                fecha
            }
        }
    `;
    const variables = { limit: 100, offset: 0 };
    const data = await fetchGraphQL(query, variables);
    const membresias = data?.data?.Membresias;

    if (!membresias || membresias.length === 0) {
        contenido.innerHTML = "<p>No hay membresías registradas.</p>";
        return;
    }

    // Generar HTML con tarjetas neutras en rejilla
    let html = '<div class="grid-container"><h3>Membresías Registradas</h3>';

    membresias.forEach(m => {
        html += `<div class="grid-item">
                    <p><strong>ID Membresía:</strong> ${m.id_membresia}</p>
                    <p><strong>Cliente:</strong> ${m.cliente?.nombre || 'N/A'}</p>
                    <p><strong>ID Tarjeta:</strong> ${m.id_tarjeta}</p>
                    <p><strong>Fecha:</strong> ${m.fecha}</p>
                 </div>`;
    });

    html += '</div>';

    contenido.innerHTML = html;
});
}