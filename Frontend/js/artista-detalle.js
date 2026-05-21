// codigo de la pagina de detalle de artista

// Llamar a la función para cargar el detalle del artista al cargar la página
document.addEventListener('DOMContentLoaded', cargarDetalleArtista);

async function cargarDetalleArtista() {

    const params = new URLSearchParams(window.location.search);
    const artistaId = params.get("id");

    if (!artistaId) {
        document.getElementById("obrasContainer").innerHTML = "<p>Artista no especificado.</p>";
    } else {
        const queryArtista = `
        query FindArtista($id: ID!) {
            findArtista(id: $id) {
                id
                nombre
                fecha_nacimiento
                nacionalidad
                biografia
                foto
            }
        }
    `;

        fetch("http://localhost:8080/query", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({
                query: queryArtista,
                variables: {id: artistaId}
            })
        })
            .then(res => res.json())
            .then(data => {
                const artista = data.data.findArtista[0];
                if (!artista) return;

                document.getElementById("artistaNombre").textContent = artista.nombre;
                document.getElementById("artistaFecha").textContent = artista.fecha_nacimiento;
                document.getElementById("artistaNacionalidad").textContent = artista.nacionalidad;
                document.getElementById("artistaBio").textContent = artista.biografia;

                let imagen;
                if (artista.foto) {
                    imagen = artista.foto
                    document.getElementById("artistaFoto").src = imagen;
                }
            });

        const queryObras = `
        query {
            Obras(limit: 100, offset: 0) {
                id
                nombre
                foto
                artista {
                    id
                    nombre
                }
            }
        }
    `;

        fetch("http://localhost:8080/query", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({query: queryObras})
        })
            .then(res => res.json())
            .then(data => {
                const obras = data.data.Obras.filter(o => o.artista.id === artistaId);
                const cont = document.getElementById("obrasContainer");

                if (obras.length === 0) {
                    cont.innerHTML = "<p>No hay obras para este artista.</p>";
                    return;
                }

                obras.forEach(o => {
                    const div = document.createElement("div");
                    div.className = "obra-card";
                    div.innerHTML = `
                <a href="obra.html?id=${o.id}">
                    <img src="${o.foto}" alt="${o.nombre}">
                    <h4>${o.nombre}</h4>
                </a>
            `;
                    cont.appendChild(div);
                });
            });
    }
}