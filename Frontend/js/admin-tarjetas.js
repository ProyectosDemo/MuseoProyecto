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

const contenido = document.getElementById("contenidoTarjetas");
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
        query ObtenerTarjetas($limit: Int, $offset: Int) {
            TarjetasCliente(limit: $limit, offset: $offset) {
                id_tarjeta
                cliente {
                    nombre
                }
                numero_tarjeta
                tipo
            }
        }
    `;

    const variables = { limit: 100, offset: 0 };
    const data = await fetchGraphQL(query, variables);
    const tarjetas = data?.data?.TarjetasCliente;

    if (!tarjetas || tarjetas.length === 0) {
        contenido.innerHTML = "<p>No hay tarjetas registradas.</p>";
        return;
    }

    // Generar HTML con tarjetas neutras en rejilla
    let html = '<div class="grid-container"><h3>Tarjetas Registradas</h3>';

    tarjetas.forEach(t => {
        html += `<div class="grid-item">
                    <p><strong>ID Tarjeta:</strong> ${t.id_tarjeta}</p>
                    <p><strong>Cliente:</strong> ${t.cliente?.nombre || 'N/A'}</p>
                    <p><strong>Número:</strong> ${t.numero_tarjeta}</p>
                    <p><strong>Tipo:</strong> ${t.tipo}</p>
                 </div>`;
    });

    html += '</div>';
    contenido.innerHTML = html;
});

}