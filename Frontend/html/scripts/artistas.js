// codigo de la pagina de artistas
//
// Llamar a la función para cargar las obras al cargar la página
document.addEventListener('DOMContentLoaded', cargarArtistas);

async function cargarArtistas() {
    // Contenedor de artistas
    const container = document.getElementById("artistasContainer");

    const query = `
                query {
                    Artistas(limit: 100, offset: 0) {
                        id
                        nombre
                        fecha_nacimiento
                        nacionalidad
                        biografia
                        foto
                    }
                }
            `;

    // Fetch GraphQL
    fetch("http://localhost:8080/query", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query })
    })
        .then(res => res.json())
        .then(data => {
            if (!data.data || !data.data.Artistas) {
                console.error("No se encontraron artistas");
                return;
            }

            const artistas = data.data.Artistas;

            artistas.forEach(artista => {
                const card = document.createElement("div");
                card.classList.add("artista-card");

                card.innerHTML = `
                        <img src="static/html/${artista.foto}" alt="${artista.nombre}">
                        <div class="artista-info">
                            <h2>${artista.nombre}</h2>
                            <p><strong>Fecha de nacimiento:</strong> ${artista.fecha_nacimiento}</p>
                            <p><strong>Nacionalidad:</strong> ${artista.nacionalidad}</p>
                            <p>${artista.biografia}</p>
                            <a href="artista-detalle.html?id=${artista.id}">Ver más</a>
                        </div>
                    `;

                container.appendChild(card);
            });
        })
        .catch(err => {
            console.error("Error cargando artistas:", err);
        });
}
