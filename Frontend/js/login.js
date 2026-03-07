// logica de la pagina de login

// Llamar a la función para manejar el login al cargar la página
 document.addEventListener('DOMContentLoaded', manejarLogin);

async function manejarLogin() {
    const form = document.getElementById("loginForm");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const login = document.getElementById("login").value.trim();
        const password = document.getElementById("password").value.trim();

        const query = `
        query {
            loginCliente(login: "${login}", password: "${password}") {
                success
                id
                nombre
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

            if (data.data.loginCliente.success) {
                localStorage.setItem("clienteNombre", data.data.loginCliente.nombre);
                localStorage.setItem("clienteId", data.data.loginCliente.id);
                window.location.href = "index.html";
            } else {
                alert("Usuario o contraseña incorrectos");
            }
        } catch (err) {
            console.error(err);
            alert("Error conectando al servidor");
        }
    });
}