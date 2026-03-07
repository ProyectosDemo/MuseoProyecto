// logica de la pagina login de trabajadores

// Llamar a la función para manejar el login al cargar la página
document.addEventListener('DOMContentLoaded', manejarLoginTrabajadores);

async function manejarLoginTrabajadores() {
    const form = document.getElementById("loginTrabajadorForm");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const login = document.getElementById("login_fake").value;
        const password = document.getElementById("password_fake").value;

        const query = `
        query {
            loginTrabajador(login: "${login}", password: "${password}") {
                success
                id
                nombre
                admin
            }
        }
    `;

        try {
            const res = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify({query})
            });

            const data = await res.json();
            const trabajador = data.data.loginTrabajador;

            if (trabajador.success) {
                localStorage.setItem("trabajadorNombre", trabajador.nombre);
                localStorage.setItem("trabajadorId", trabajador.id);
                localStorage.setItem("trabajadorAdmin", trabajador.admin);

                if (trabajador.admin) {
                    window.location.href = "admin.html";
                } else {
                    window.location.href = "index.html";
                }
            } else {
                alert("Usuario o contraseña incorrectos");
            }
        } catch (err) {
            console.error(err);
            alert("Error conectando al servidor");
        }
    });
}