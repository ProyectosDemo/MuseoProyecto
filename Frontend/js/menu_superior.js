// vas a usar esto en todas las paginas
document.addEventListener('DOMContentLoaded', header);

function header() {

    // busca el header existente para reemplazarlo
    let mainHeader = document.getElementById('header');
    if (mainHeader) {
        mainHeader.replaceWith(crearHeader());
    } else {
        // Si no existe, lo insertamos al inicio del body
        document.body.insertBefore(crearHeader(), document.body.firstChild);
    }

    // Lógica del menú superior (sesión, botones dinámicos)
    const navMenu = document.getElementById("navMenu");
    const loginButton = document.getElementById("loginButton");

    const nombreCliente = localStorage.getItem("clienteNombre");
    const nombreTrabajador = localStorage.getItem("trabajadorNombre");
    const esAdmin = localStorage.getItem("trabajadorAdmin") === "true";
    let userType = "";

    if (nombreCliente) {
        loginButton.textContent = `Bienvenido, ${nombreCliente}`;
        loginButton.href = "#";
        userType = "cliente";
    } else if (nombreTrabajador) {
        loginButton.textContent = `Bienvenido, ${nombreTrabajador}`;
        loginButton.href = "#";
        userType = esAdmin ? "admin" : "trabajador";
    }

    if (userType) {
        navMenu.insertBefore(loginButton, navMenu.firstChild);

        if (userType === "trabajador" || esAdmin) {
            const reservasBtn = document.createElement("a");
            reservasBtn.href = "reservas.html";
            reservasBtn.textContent = "Reservas";
            navMenu.appendChild(reservasBtn);
        }

        if (esAdmin) {
            const adminBtn = document.createElement("a");
            adminBtn.href = "admin.html";
            adminBtn.textContent = "Administración";
            navMenu.appendChild(adminBtn);
        }

        if (userType === "cliente") {
            const cuentaBtn = document.createElement("a");
            cuentaBtn.href = "cuenta-cliente.html";
            cuentaBtn.textContent = "Mi Cuenta";
            cuentaBtn.style.marginLeft = "10px";
            navMenu.appendChild(cuentaBtn);
        }

        const logoutBtn = document.createElement("a");
        logoutBtn.href = "#";
        logoutBtn.textContent = "Cerrar Sesión";
        logoutBtn.addEventListener("click", () => {
            if (userType === "cliente") {
                localStorage.removeItem("clienteNombre");
                localStorage.removeItem("clienteId");
            } else {
                localStorage.removeItem("trabajadorNombre");
                localStorage.removeItem("trabajadorId");
                localStorage.removeItem("trabajadorAdmin");
            }
            location.reload();
        });
        navMenu.appendChild(logoutBtn);
    }
}

function crearHeader() {
    const header = document.createElement('header');
    header.innerHTML = `
        <h1>Museo de Arte Contemporáneo</h1>
        <nav id="navMenu">
            <a href="index.html">Inicio</a>
            <a href="exposiciones.html">Exposiciones</a>
            <a href="artistas.html">Artistas</a>
            <a href="login.html" id="loginButton">Iniciar Sesión</a>
        </nav>
    `;
    return header;
}