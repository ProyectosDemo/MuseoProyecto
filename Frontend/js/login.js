// logica de la pagina de login

// Llamar a la función para manejar el login al cargar la página
 document.addEventListener('DOMContentLoaded', manejarLogin);

function manejarLogin() {

    const form = document.getElementById("loginForm");
    const loginInput = document.getElementById("login");
    const passwordInput = document.getElementById("password");
    const errorLogin = document.getElementById("errorLogin");
    const errorPassword = document.getElementById("errorPassword");

    form.addEventListener("submit", async (e) => {
        e.preventDefault(); // evita que recargue la pagina al enviar el formulario

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