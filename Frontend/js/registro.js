// logica del registro de clientes

// Llamar a la función para manejar el registro al cargar la página
document.addEventListener('DOMContentLoaded', manejarRegistro);

const opcionesPreguntas = [
    "Nombre de tu primera mascota?",
    "Ciudad donde naciste?",
    "Tu comida favorita?",
    "Nombre de tu mejor amigo de infancia?",
    "Tu película favorita?",
    "Primer colegio al que fuiste?"
];

const selectsPreguntas = [
    document.getElementById("pregunta1"),
    document.getElementById("pregunta2"),
    document.getElementById("pregunta3")
];

function actualizarSelects() {
    const seleccionadas = selectsPreguntas.map(s => s.value).filter(v => v !== "");
    selectsPreguntas.forEach(select => {
        const valorActual = select.value;
        select.innerHTML = `<option value="" disabled ${valorActual === "" ? "selected" : ""}>Selecciona una pregunta</option>`;
        opcionesPreguntas.forEach(op => {
            if (!seleccionadas.includes(op) || op === valorActual) {
                const opt = document.createElement("option");
                opt.value = op;
                opt.textContent = op;
                if (op === valorActual) opt.selected = true;
                select.appendChild(opt);
            }
        });
    });
}

async function manejarRegistro() {
    actualizarSelects();

    selectsPreguntas.forEach(select => select.addEventListener("change", actualizarSelects));

    const form = document.getElementById("registroForm");

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const nombre = document.getElementById("nombre").value.trim();
        const email = document.getElementById("email").value.trim();
        const telefono = document.getElementById("telefono").value.trim();
        const login = document.getElementById("login").value.trim();
        const password = document.getElementById("password").value.trim();
        const codigoSeguridadCliente = generarCodigoSeguridad(8);

        const numeroTarjeta = document.getElementById("numero_tarjeta").value.trim();
        const tipo = document.getElementById("tipo").value;

        try {
            // cliente
            const createClienteQuery = `
            mutation CreateCliente($input: NewCliente!) {
                createCliente(input: $input) {
                    id
                    nombre
                }
            }
        `;
            const clienteVariables = { input: { nombre, email, telefono, login, password, codigo_seguridad: codigoSeguridadCliente } };
            const resCliente = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: createClienteQuery, variables: clienteVariables })
            });
            const dataCliente = await resCliente.json();
            if (dataCliente.errors) throw new Error(dataCliente.errors[0].message);
            const clienteId = dataCliente.data.createCliente?.id;
            if (!clienteId) throw new Error("No se pudo obtener el ID del cliente.");

            //tarjeta
            const createTarjetaQuery = `
            mutation CreateTarjeta($input: NewTarjetaCliente!) {
                createTarjetaCliente(input: $input) {
                    id_tarjeta
                }
            }
        `;
            const tarjetaVariables = { input: { id_cliente: clienteId, numero_tarjeta: numeroTarjeta, tipo } };
            const resTarjeta = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: createTarjetaQuery, variables: tarjetaVariables })
            });
            const dataTarjeta = await resTarjeta.json();
            if (dataTarjeta.errors) throw new Error(dataTarjeta.errors[0].message);
            const tarjetaId = dataTarjeta.data.createTarjetaCliente?.id_tarjeta;
            if (!tarjetaId) throw new Error("No se pudo obtener el ID de la tarjeta.");

            // membresia
            const fechaHoy = new Date().toISOString().split("T")[0]; // YYYY-MM-DD
            const createMembresiaQuery = `
            mutation CreateMembresia($input: NewMembresia!) {
                createMembresia(input: $input) {
                    id_membresia
                    id_cliente
                    id_tarjeta
                    fecha
                }
            }
        `;
            const membresiaVariables = { input: { id_cliente: clienteId, id_tarjeta: tarjetaId, fecha: fechaHoy } };
            const resMembresia = await fetch("http://localhost:8080/query", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ query: createMembresiaQuery, variables: membresiaVariables })
            });
            const dataMembresia = await resMembresia.json();
            if (dataMembresia.errors) throw new Error(dataMembresia.errors[0].message);
            const membresiaId = dataMembresia.data.createMembresia?.id_membresia;
            if (!membresiaId) throw new Error("No se pudo crear la membresía.");

            // preguntas
            for (let i = 1; i <= 3; i++) {
                const createPreguntaQuery = `
                mutation CreatePreguntas($input: NewPreguntas!) {
                    createPreguntas(input: $input) {
                        id
                        pregunta
                        respuesta
                    }
                }
            `;
                const preguntaVariables = {
                    input: {
                        id_cliente: clienteId,
                        pregunta: document.getElementById(`pregunta${i}`).value,
                        respuesta: document.getElementById(`respuesta${i}`).value.trim()
                    }
                };
                const resPregunta = await fetch("http://localhost:8080/query", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ query: createPreguntaQuery, variables: preguntaVariables })
                });
                const dataPregunta = await resPregunta.json();
                if (dataPregunta.errors) throw new Error(dataPregunta.errors[0].message);
            }

            alert(`Cuenta creada correctamente.\nCódigo de seguridad: ${codigoSeguridadCliente}\nID Membresía: ${membresiaId}`);
            window.location.href = "login.html";

        } catch (err) {
            console.error(err);
            alert("Error al crear la cuenta: " + err.message);
        }
    });
}

function generarCodigoSeguridad(longitud = 8) {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    for (let i = 0; i < longitud; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}