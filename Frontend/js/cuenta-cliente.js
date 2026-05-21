// logica de la pagina cuenta del cliente

document.addEventListener('DOMContentLoaded', manejarCuentaCliente);

let clienteNombre = localStorage.getItem("clienteNombre") || null;
let clienteID = localStorage.getItem("clienteId") || null;
const navMenu = document.getElementById("navMenu");
const loginButton = document.getElementById("loginButton");

const loginForm = document.getElementById('login-form');
const editForm = document.getElementById('edit-form');
const loginSection = document.getElementById('login-section');
const editSection = document.getElementById('edit-section');
const loginError = document.getElementById('login-error');
const editSuccess = document.getElementById('edit-success');
const editError = document.getElementById('edit-error');
const recuperarBtn = document.getElementById('recuperar-codigo');
const codigoDisplay = document.getElementById('codigo-display');
const preguntasSection = document.getElementById('preguntas-section');
const preguntasContainer = document.getElementById('preguntas-container');
const preguntasForm = document.getElementById('preguntas-form');
const preguntasError = document.getElementById('preguntas-error');

async function manejarCuentaCliente() {

    renderNav();
    // VERIFICACIÓN DE CONTRASEÑA
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        loginError.textContent = '';

        const login = document.getElementById('login').value;
        const passwordActual = document.getElementById('password').value;

        const query = `
        query {
            loginCliente(login: "${login}", password: "${passwordActual}") {
                success
                id
                nombre
            }
        }
    `;

        try {
            const res = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query })
            });
            const data = await res.json();

            if (data.data && data.data.loginCliente && data.data.loginCliente.success) {
                clienteID = data.data.loginCliente.id;
                clienteNombre = data.data.loginCliente.nombre;
                localStorage.setItem("clienteId", clienteID);
                localStorage.setItem("clienteNombre", clienteNombre);
                renderNav();

                const clienteQuery = `
                query {
                    findCliente(id: "${clienteID}") {
                        id
                        nombre
                        email
                        telefono
                    }
                }
            `;
                const clienteRes = await fetch("http://localhost:8080/query", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ query: clienteQuery })
                });
                const clienteData = await clienteRes.json();
                if (clienteData.data && clienteData.data.findCliente) {
                    const cliente = clienteData.data.findCliente[0];
                    document.getElementById('nombre').value = cliente.nombre;
                    document.getElementById('email').value = cliente.email;
                    document.getElementById('telefono').value = cliente.telefono;

                    loginSection.style.display = 'none';
                    editSection.style.display = 'block';
                }
            } else {
                loginError.textContent = 'Contraseña incorrecta';
            }
        } catch (err) {
            loginError.textContent = 'Error al verificar contraseña';
            console.error(err);
        }
    });

    editForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        editError.textContent = '';
        editSuccess.textContent = '';

        const nombre = document.getElementById('nombre').value;
        const email = document.getElementById('email').value;
        const telefono = document.getElementById('telefono').value;
        const password = document.getElementById('password-new').value;

        const mutation = `
        mutation {
            updateCliente(input: {
                id: "${clienteID}",
                nombre: "${nombre}",
                email: "${email}",
                telefono: "${telefono}",
                ${password ? `password: "${password}",` : ''}
            }) {
                id
                nombre
            }
        }
    `;

        try {
            const res = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: mutation })
            });
            const data = await res.json();
            if (data.data && data.data.updateCliente) {
                editSuccess.textContent = 'Datos actualizados correctamente';
                if (password) document.getElementById('password-new').value = '';
            } else {
                editError.textContent = 'Error al actualizar los datos';
            }
        } catch (err) {
            editError.textContent = 'Error al actualizar los datos';
            console.error(err);
        }
    });

    recuperarBtn.addEventListener('click', async () => {
        preguntasError.textContent = '';
        codigoDisplay.textContent = '';
        preguntasContainer.innerHTML = '';
        preguntasSection.style.display = 'block';

        if (!clienteID) {
            preguntasError.textContent = 'No se encontró el cliente. Debes iniciar sesión.';
            return;
        }

        try {
            const query = `
            query {
                findPreguntasByCliente(id_cliente: "${clienteID}") {
                    pregunta
                    pregunta
                    respuesta
                }
            }
        `;
            const res = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query })
            });
            const data = await res.json();

            if (!data.data || !data.data.findPreguntasByCliente || data.data.findPreguntasByCliente.length === 0) {
                preguntasError.textContent = 'No hay preguntas registradas para este cliente.';
                console.error(data.errors || data);
                return;
            }

            const preguntas = data.data.findPreguntasByCliente;

            preguntas.forEach((p, index) => {
                const div = document.createElement('div');
                div.innerHTML = `
                <label>${p.pregunta}</label>
                <input type="text" name="respuesta" data-respuesta="${p.respuesta}" required>
            `;
                preguntasContainer.appendChild(div);
            });
        } catch (err) {
            preguntasError.textContent = 'Error al cargar preguntas';
            console.error(err);
        }
    });

    preguntasForm.addEventListener('submit', (e) => {
        e.preventDefault();
        preguntasError.textContent = '';

        const inputs = preguntasForm.querySelectorAll('input[name="respuesta"]');
        let allCorrect = true;

        inputs.forEach(input => {
            const correcta = input.dataset.respuesta.trim();
            const respuestaUsuario = input.value.trim();
            if (respuestaUsuario !== correcta) allCorrect = false;
        });

        if (allCorrect) {
            // Mostrar codigo de seguridad
            fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    query: `
                    query {
                        findCliente(id: "${clienteID}") {
                            codigo_seguridad
                        }
                    }
                `
                })
            })
                .then(res => res.json())
                .then(data => {
                    if (data.data && data.data.findCliente) {
                        codigoDisplay.textContent = `Código de seguridad: ${data.data.findCliente[0].codigo_seguridad}`;
                    } else {
                        preguntasError.textContent = 'No se pudo recuperar el código.';
                    }
                })
                .catch(err => {
                    preguntasError.textContent = 'Error al recuperar el código.';
                    console.error(err);
                });
        } else {
            preguntasError.textContent = 'Alguna respuesta es incorrecta';
        }
    });
}


function renderNav() {
    navMenu.innerHTML = '';
    const homeLink = document.createElement('a');
    homeLink.href = "index.html";
    homeLink.textContent = "Inicio";
    navMenu.appendChild(homeLink);

    if (clienteNombre) {
        loginButton.textContent = `Bienvenido, ${clienteNombre}`;
        loginButton.href = "#";
        navMenu.appendChild(loginButton);

        const cuentaBtn = document.createElement("a");
        cuentaBtn.href = "cuenta-cliente.html";
        cuentaBtn.textContent = "Mi Cuenta";
        cuentaBtn.style.marginLeft = "10px";
        navMenu.appendChild(cuentaBtn);

        const logoutBtn = document.createElement("a");
        logoutBtn.href = "#";
        logoutBtn.textContent = "Cerrar Sesión";
        logoutBtn.style.marginLeft = "10px";
        logoutBtn.addEventListener("click", () => {
            localStorage.removeItem("clienteNombre");
            localStorage.removeItem("clienteId");
            location.reload();
        });
        navMenu.appendChild(logoutBtn);
    } else {
        loginButton.textContent = "Iniciar Sesión";
        loginButton.href = "login.html";
        navMenu.appendChild(loginButton);
    }
}
