document.addEventListener('DOMContentLoaded', manejarArtistas);

async function manejarArtistas() {
    await cargarGeneros();
    await cargarArtistas();
}

async function cargarGeneros() {
    const filtroGeneroContainer = document.getElementById("filtroGenero");
    const query = `
        query {
            Genero(limit: 100, offset: 0) {
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
        const generos = data.data?.Genero || [];

        generos.forEach(g => {
            const option = document.createElement("option");
            option.value = g.id;
            option.textContent = g.nombre;
            filtroGeneroContainer.appendChild(option);
        });

        filtroGeneroContainer.addEventListener("change", () => {
            cargarArtistas(filtroGeneroContainer.value);
        });

    } catch (err) {
        console.error("Error cargando géneros:", err);
    }
}

async function cargarArtistas(idGenero = "") {
    const gallery = document.getElementById("gallery");
    gallery.innerHTML = "";

    let query, variables;
    if (idGenero) {
        query = `
            query($idGenero: ID!) {
                FindArtistaGeneroByGenero(id_genero: $idGenero) {
                    artista {
                        id
                        nombre
                        fecha_nacimiento
                        nacionalidad
                        biografia
                        foto
                    }
                }
            }
        `;
        variables = { idGenero };
    } else {
        query = `
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
        variables = null;
    }

    try {
        const res = await fetch("http://localhost:8080/query", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ query, variables })
        });
        const data = await res.json();

        let artistas;
        if (idGenero) {
            artistas = data.data?.FindArtistaGeneroByGenero?.map(ag => ag.artista) || [];
        } else {
            artistas = data.data?.Artistas || [];
        }

        if (!artistas.length) {
            gallery.innerHTML = "<p>No hay artistas disponibles</p>";
            return;
        }

        artistas.forEach(artista => {
            const artCard = document.createElement("div");
            artCard.className = "art-card";
            artCard.innerHTML = `
                <div class="art-info">
                    <a href="artista-detalle.html?id=${artista.id}">
                        <img src="${artista.foto}" alt="${artista.nombre}">
                        <h4>${artista.nombre}</h4>
                    </a>
                    <p>${artista.fecha_nacimiento}</p>
                    <p>${artista.nacionalidad}</p>
                </div>
            `;
            gallery.appendChild(artCard);
        });

    } catch (err) {
        console.error(err);
        gallery.innerHTML = "<p>Error cargando artistas</p>";
    }
}