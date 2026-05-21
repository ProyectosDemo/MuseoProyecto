// logica de la pagina login de trabajadores

// Llamar a la función para manejar el login al cargar la página
document.addEventListener('DOMContentLoaded', manejarLoginTrabajadores);

async function manejarLoginTrabajadores() {
    const form = document.getElementById("loginTrabajadorForm");
    const loginInput = document.getElementById("login_fake");
    const passwordInput = document.getElementById("password_fake");
    const errorLogin = document.getElementById("errorLogin");
    const errorPassword = document.getElementById("errorPassword");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        // limpia errores de antes
        errorLogin.textContent = "";
        errorPassword.textContent = "";
        // saca los valores de los campos
        const login = loginInput.value.trim();
        const password = passwordInput.value.trim();
        // flag de errores
        let errores = false;

        if (login === "") {
            if (errorLogin) {
                errorLogin.textContent = "Debes llenar el campo usuario.";
            }
            loginInput.classList.add("input-error");
            errores = true;
        }

        if (password === "") {
            if (errorPassword) {
                errorPassword.textContent = "Debes llenar el campo contraseña.";
            }
            passwordInput.classList.add("input-error");
            errores = true;
        }

        // si no ingreso datos no hacemos la consulta
        if (errores) {
            // Enfocar el primer campo con error
            if (login === "") {
                loginInput.focus();
            } else if (password === "") {
                passwordInput.focus();
            }
            return;
        }

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

                if (trabajador.admin) { // lo mandamos a la pagina de administracion
                    window.location.href = "admin.html";
                } else {
                    window.location.href = "index.html";
                }

            } else {
                // Error de credenciales
                if (errorPassword) {
                    errorPassword.textContent = "Usuario o contraseña incorrectos.";
                }
            }
        } catch (err) {
            console.error(err);
            alert("Error conectando al servidor");
        }
    });
}