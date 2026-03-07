document.addEventListener('DOMContentLoaded', manejarExposiciones);

async function manejarExposiciones() {
    await cargarTodasObras();
    await cargarFiltros();
}

let todasObras = [];

// Cargar todas las obras
async function cargarTodasObras() {
    const query = `
        query {
            Obras(limit: 100, offset: 0) {
                id
                nombre
                foto
                artista {
                    id
                    nombre
                }
                genero {
                    id
                    nombre
                }
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
        todasObras = data.data.Obras || [];
        mostrarObras(todasObras);
    } catch (err) {
        console.error(err);
        document.getElementById("gallery").innerHTML = "<p>Error cargando obras</p>";
    }
}

// Mostrar las obras en el DOM
function mostrarObras(obras) {
    const gallery = document.getElementById("gallery");
    gallery.innerHTML = "";

    if (!obras || obras.length === 0) {
        gallery.innerHTML = "<p>No hay obras disponibles</p>";
        return;
    }

    obras.forEach(ob => {
        const artCard = document.createElement("div");
        artCard.className = "art-card";

        artCard.innerHTML = `
            <div class="art-info">
                <a href="obra.html?id=${ob.id}">
                    <img src="${ob.foto}" alt="${ob.nombre}">
                    <h4>${ob.nombre}</h4>
                </a>
                <p>${ob.artista ? ob.artista.nombre : ""}</p>
                <p>${ob.genero ? ob.genero.nombre : ""}</p>
            </div>
        `;
        gallery.appendChild(artCard);
    });
}

// Cargar filtros de artista y género
async function cargarFiltros() {
    const artistasQuery = `
        query { Artistas(limit:100, offset:0) { id nombre } }
    `;
    const generosQuery = `
        query { Genero(limit:100, offset:0) { id nombre } }
    `;

    const [artistasRes, generosRes] = await Promise.all([
        fetch("http://localhost:8080/query", { method:"POST", headers:{"Content-Type":"application/json"}, body:JSON.stringify({query: artistasQuery}) }).then(r => r.json()),
        fetch("http://localhost:8080/query", { method:"POST", headers:{"Content-Type":"application/json"}, body:JSON.stringify({query: generosQuery}) }).then(r => r.json())
    ]);

    const artistaSelect = document.getElementById("filtroArtista");
    const generoSelect = document.getElementById("filtroGenero");

    if (artistaSelect) {
        artistaSelect.innerHTML = '<option value="">Todos los artistas</option>';
        (artistasRes.data.Artistas || []).forEach(a => {
            artistaSelect.innerHTML += `<option value="${a.id}">${a.nombre}</option>`;
        });
        artistaSelect.addEventListener('change', filtrarObras);
    }

    if (generoSelect) {
        generoSelect.innerHTML = '<option value="">Todos los géneros</option>';
        (generosRes.data.Genero || []).forEach(g => {
            generoSelect.innerHTML += `<option value="${g.id}">${g.nombre}</option>`;
        });
        generoSelect.addEventListener('change', filtrarObras);
    }
}

// Filtrar obras por artista y/o género
function filtrarObras() {
    const idArtista = document.getElementById("filtroArtista")?.value;
    const idGenero = document.getElementById("filtroGenero")?.value;

    let obrasFiltradas = [...todasObras];

    if (idArtista) obrasFiltradas = obrasFiltradas.filter(o => String(o.artista.id) === String(idArtista));
    if (idGenero) obrasFiltradas = obrasFiltradas.filter(o => String(o.genero.id) === String(idGenero));

    mostrarObras(obrasFiltradas);
}