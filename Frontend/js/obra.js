document.addEventListener('DOMContentLoaded', cargarObra);

async function cargarObra() {
    const params = new URLSearchParams(window.location.search);
    const obraId = params.get("id");
    if (!obraId) return;

    try {
        const query = `
            query FindObraEscultura($id_obra: ID!) {
                findObra(id: $id_obra) {
                    id
                    nombre
                    artista { id nombre }
                    genero { nombre }
                    precio
                    fecha_creacion
                    status
                    foto
                }
                findEscultura(id_obra: $id_obra) {
                    id_obra
                    material
                    peso
                    dimensiones
                }
            }
        `;

        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query, variables: { id_obra: String(obraId) } })
        });

        const data = await res.json();
        const obra = data?.data?.findObra?.[0];
        const escultura = data?.data?.findEscultura;

        if (!obra) return;

        document.getElementById("obraFoto").src = obra.foto;
        document.getElementById("obraNombre").textContent = obra.nombre;
        document.getElementById("obraArtista").innerHTML = obra.artista ? `<a class="link-link" href="artista-detalle.html?id=${obra.artista?.id}">${obra.artista?.nombre}</a>` : "";
        document.getElementById("obraGenero").textContent = obra.genero?.nombre || "";
        document.getElementById("obraPrecio").textContent = "$" + obra.precio;
        document.getElementById("obraFecha").textContent = obra.fecha_creacion;
        document.getElementById("obraEstatus").textContent = obra.status;

        const divMaterial = document.getElementById("divMaterial");
        const divPeso = document.getElementById("divPeso");
        const divDimensiones = document.getElementById("divDimensiones");

        if (escultura) {
            divMaterial.style.display = "flex";
            divPeso.style.display = "flex";
            divDimensiones.style.display = "flex";

            document.getElementById("obraMaterial").textContent = escultura.material;
            document.getElementById("obraPeso").textContent = escultura.peso + " kg";
            document.getElementById("obraDimensiones").textContent = escultura.dimensiones;
        } else {
            divMaterial.style.display = "none";
            divPeso.style.display = "none";
            divDimensiones.style.display = "none";
        }

        const clienteId = localStorage.getItem("clienteId");
        if (obra.status === "DISPONIBLE" && clienteId) {
            // Solo crear el botón si no existe
            if (!document.getElementById("btnComprar")) {
                const btnComprar = document.createElement("button");
                btnComprar.id = "btnComprar";
                btnComprar.textContent = "Comprar";

                btnComprar.addEventListener("click", () => {
                    // Solo crear el formulario si no existe
                    if (!document.getElementById("formCodigoSeguridad")) {
                        const formDiv = document.createElement("div");
                        formDiv.id = "formCodigoSeguridad";
                        formDiv.innerHTML = `
                            <input type="password" id="codigoSeguridad" placeholder="Código de seguridad">
                            <button id="btnConfirmar">Confirmar</button>
                        `;
                        document.querySelector(".foto-boton").appendChild(formDiv);

                        formDiv.querySelector("#btnConfirmar").addEventListener("click", async () => {
                            const codigoIngresado = formDiv.querySelector("#codigoSeguridad").value;
                            if (!codigoIngresado) { alert("Debe ingresar el código de seguridad."); return; }

                            try {
                                const queryCliente = `
                                    query FindCliente($id: ID!) { findCliente(id: $id) { id codigo_seguridad } }
                                `;
                                const resCliente = await fetch("http://localhost:8080/query", {
                                    method: "POST",
                                    headers: { "Content-Type": "application/json" },
                                    body: JSON.stringify({ query: queryCliente, variables: { id: clienteId } })
                                });
                                const dataCliente = await resCliente.json();
                                const cliente = dataCliente?.data?.findCliente?.[0];

                                if (!cliente) throw new Error("Cliente no encontrado");
                                if (cliente.codigo_seguridad !== codigoIngresado) {
                                    alert("Código de seguridad incorrecto.");
                                    throw new Error("Código incorrecto");
                                }

                                const mutationUpdateObra = `
                                    mutation UpdateObra($id: ID!, $status: StatusObra!) {
                                        updateObra(input: { id: $id, status: $status }) { id status }
                                    }
                                `;
                                await fetch("http://localhost:8080/query", {
                                    method: "POST",
                                    headers: { "Content-Type": "application/json" },
                                    body: JSON.stringify({ query: mutationUpdateObra, variables: { id: obra.id, status: "RESERVADA" } })
                                });

                                const mutationCreateOrden = `
                                    mutation CreateOrden($input: NewOrden!) {
                                        createOrden(input: $input) { id status }
                                    }
                                `;
                                const fechaActual = new Date().toISOString();
                                await fetch("http://localhost:8080/query", {
                                    method: "POST",
                                    headers: { "Content-Type": "application/json" },
                                    body: JSON.stringify({
                                        query: mutationCreateOrden,
                                        variables: { input: { id_obra: obra.id, id_cliente: clienteId, id_trabajador: null, fecha: fechaActual, status: "PENDIENTE" } }
                                    })
                                });

                                alert("Obra reservada y orden creada exitosamente!");
                                document.getElementById("obraEstatus").textContent = "RESERVADA";
                                formDiv.remove();
                                btnComprar.remove();
                            } catch (err) {
                                console.error(err);
                            }
                        });
                    }
                });

                document.querySelector(".foto-boton").appendChild(btnComprar);
            }
        }

    } catch (err) {
        console.error(err);
    }
}