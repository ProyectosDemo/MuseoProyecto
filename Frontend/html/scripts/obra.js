// logica de la pagina obra

// Llamar a la función para cargar las obras al cargar la página
document.addEventListener('DOMContentLoaded', cargarObra);

async function cargarObra() {
    const params = new URLSearchParams(window.location.search);
    const obraId = params.get("id");

    if (obraId) {
        const query = `
        query FindObra($id: ID!) {
            findObra(id: $id) {
                id
                nombre
                artista { nombre }
                genero { nombre }
                precio
                fecha_creacion
                status
                foto
                material
                peso
                dimensiones
            }
        }
    `;

        fetch("http://localhost:8080/query", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({query, variables: {id: String(obraId)}})
        })
            .then(res => res.json())
            .then(data => {
                const obra = data?.data?.findObra?.[0];
                if (!obra) return;

                document.getElementById("obraFoto").src = obra.foto;
                document.getElementById("obraNombre").textContent = obra.nombre;
                document.getElementById("obraArtista").textContent = obra.artista?.nombre || "";
                document.getElementById("obraGenero").textContent = obra.genero?.nombre || "";
                document.getElementById("obraPrecio").textContent = "$" + obra.precio;
                document.getElementById("obraFecha").textContent = obra.fecha_creacion;
                document.getElementById("obraEstatus").textContent = obra.status;
                document.getElementById("obraMaterial").textContent = obra.material;
                document.getElementById("obraPeso").textContent = obra.peso + " kg";
                document.getElementById("obraDimensiones").textContent = obra.dimensiones;

                // Botón comprar solo si DISPONIBLE y hay cliente logueado
                if (obra.status === "DISPONIBLE" && clienteId) {
                    const btnComprar = document.createElement("button");
                    btnComprar.id = "btnComprar";
                    btnComprar.textContent = "Comprar";

                    btnComprar.addEventListener("click", () => {
                        const formDiv = document.createElement("div");
                        formDiv.id = "formCodigoSeguridad";
                        formDiv.innerHTML = `
                    <input type="password" id="codigoSeguridad" placeholder="Código de seguridad">
                    <button id="btnConfirmar">Confirmar</button>
                `;
                        document.querySelector(".foto-boton").appendChild(formDiv);

                        formDiv.querySelector("#btnConfirmar").addEventListener("click", () => {
                            const codigoIngresado = formDiv.querySelector("#codigoSeguridad").value;
                            if (!codigoIngresado) {
                                alert("Debe ingresar el código de seguridad.");
                                return;
                            }

                            // Validar cliente
                            const queryCliente = `
                        query FindCliente($id: ID!) {
                            findCliente(id: $id) { id codigo_seguridad }
                        }
                    `;
                            fetch("http://localhost:8080/query", {
                                method: "POST",
                                headers: {"Content-Type": "application/json"},
                                body: JSON.stringify({query: queryCliente, variables: {id: clienteId}})
                            })
                                .then(res => res.json())
                                .then(data => {
                                    const cliente = data?.data?.findCliente?.[0];
                                    if (!cliente) throw new Error("Cliente no encontrado");
                                    if (cliente.codigo_seguridad !== codigoIngresado) {
                                        alert("Código de seguridad incorrecto.");
                                        throw new Error("Código incorrecto");
                                    }

                                    // Actualizar obra a RESERVADA
                                    return fetch("http://localhost:8080/query", {
                                        method: "POST",
                                        headers: {"Content-Type": "application/json"},
                                        body: JSON.stringify({
                                            query: `
                                    mutation UpdateObra($id: ID!, $status: StatusObra!) {
                                        updateObra(input: { id: $id, status: $status }) { id status }
                                    }
                                `,
                                            variables: {id: obra.id, status: "RESERVADA"}
                                        })
                                    });
                                })
                                .then(() => {
                                    // Crear orden PENDIENTE
                                    const mutationCreateOrden = `
                            mutation CreateOrden($input: NewOrden!) {
                                createOrden(input: $input) { id status }
                            }
                        `;
                                    const fechaActual = new Date().toISOString();
                                    return fetch("http://localhost:8080/query", {
                                        method: "POST",
                                        headers: {"Content-Type": "application/json"},
                                        body: JSON.stringify({
                                            query: mutationCreateOrden,
                                            variables: {
                                                input: {
                                                    id_obra: obra.id,
                                                    id_cliente: clienteId,
                                                    id_trabajador: null,
                                                    fecha: fechaActual,
                                                    status: "PENDIENTE"
                                                }
                                            }
                                        })
                                    });
                                })
                                .then(res => res.json())
                                .then(data => {
                                    alert("Obra reservada y orden creada exitosamente!");
                                    document.getElementById("obraEstatus").textContent = "RESERVADA";
                                    formDiv.remove();
                                    btnComprar.remove();
                                })
                                .catch(err => console.error(err));
                        });
                    });

                    document.querySelector(".foto-boton").appendChild(btnComprar);
                }
            })
    }
}