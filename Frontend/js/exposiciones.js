document.addEventListener('DOMContentLoaded', manejarExposiciones);

async function manejarExposiciones() {
    const catalogoInicial = await obtenerCatalogoCompleto();
    mostrarObras(catalogoInicial);
    await cargarFiltros();
}


async function obtenerCatalogoCompleto() {
    const query = `{ Obras(limit: 250, offset: 0) { id nombre foto precio status artista { id nombre } genero { id nombre } } }`;
    try {
        const res = await fetch("http://localhost:8080/query", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ query }) });
        const data = await res.json();
        return data.data.Obras || [];
    } catch (err) {
        console.error(err);
        document.getElementById("gallery").innerHTML = "<p>Error cargando obras</p>";
    }
}

async function cargarObrasPorPrecioAsc()  {
    const query = `{ ObrasPorPrecio(limit:250, offset:0) { id nombre precio foto artista { id nombre } genero { id nombre } } }`;
    try {
        const res = await fetch("http://localhost:8080/query", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ query }) });
        const data = await res.json();
        return data.data.ObrasPorPrecio || [];
    } catch (err) {
        console.error(err);
        document.getElementById("gallery").innerHTML = "<p>Error cargando obras</p>";
    }
}

async function cargarObrasPorPrecioDesc() {
    const query = `{ ObrasPorPrecioDesc(limit:250, offset:0) { id nombre precio foto artista { id nombre } genero { id nombre } } }`;
    try {
        const res = await fetch("http://localhost:8080/query", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ query }) });
        const data = await res.json();
        return data.data.ObrasPorPrecioDesc || [];
    } catch (err) {
        console.error(err);
        document.getElementById("gallery").innerHTML = "<p>Error cargando obras</p>";
    }
}

async function cargarObrasPorGenero(idGenero) {
    console.log("cargando genero"); // original
    
    const query = `{ ObrasPorGenero(idGenero: ${idGenero}, limit:250, offset:0) { id nombre precio foto artista { id nombre } genero { id nombre } } }`;
    try {
        const res = await fetch("http://localhost:8080/query", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ query }) });
        const data = await res.json();
        return data.data.ObrasPorGenero || [];
    } catch (err) {
        console.error(err);
        document.getElementById("gallery").innerHTML = "<p>Error cargando obras</p>";
    }
}

async function cargarObrasPorDisponibilidad() {
    const query = `{ ObrasPorDisponibilidad(limit:250, offset:0) { id nombre precio foto status artista { id nombre } genero { id nombre } } }`;
    try {
        const res = await fetch("http://localhost:8080/query", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ query }) });
        const data = await res.json();
        return data.data.ObrasPorDisponibilidad || [];
    } catch (err) {
        console.error(err);
        return [];
    }
}


// Cargar filtros de artista y género
async function cargarFiltros() {
    const artistasQuery = `query { Artistas(limit:250, offset:0) { id nombre } }`;
    const generosQuery = `query { Genero(limit:250, offset:0) { id nombre } }`;

    try {
        const [artistasRes, generosRes] = await Promise.all([
            fetch("http://localhost:8080/query", { method:"POST", headers:{"Content-Type":"application/json"}, body:JSON.stringify({query: artistasQuery}) }).then(r => r.json()),
            fetch("http://localhost:8080/query", { method:"POST", headers:{"Content-Type":"application/json"}, body:JSON.stringify({query: generosQuery}) }).then(r => r.json())
        ]);

        const artistaSelect = document.getElementById("filtroArtista");
        const generoSelect = document.getElementById("filtroGenero");
        const precioSelect = document.getElementById("filtroPrecio");
        const disponibilidadSelect = document.getElementById("filtroDisponibilidad");

        if (artistaSelect) {
            artistaSelect.innerHTML = '<option value="">Todos los artistas</option>';
            (artistasRes.data.Artistas || []).forEach(a => { artistaSelect.innerHTML += `<option value="${a.id}">${a.nombre}</option>`; });
            artistaSelect.addEventListener('change', filtrarObras);
        }

        if (generoSelect) {
            generoSelect.innerHTML = '<option value="">Todos los géneros</option>';
            (generosRes.data.Genero || []).forEach(g => { generoSelect.innerHTML += `<option value="${g.id}">${g.nombre}</option>`; });
            generoSelect.addEventListener('change', filtrarObras);
        }

        if (precioSelect) {
            precioSelect.addEventListener('change', filtrarObras);
        }
        if (disponibilidadSelect) {
            disponibilidadSelect.addEventListener('change', filtrarObras);
        }
    } catch (error) {
        console.error(error);
    }
}

async function filtrarObras(){
    const idArtista = document.getElementById("filtroArtista")?.value;
    const idGenero = document.getElementById("filtroGenero")?.value;
    const ordenPrecio = document.getElementById("filtroPrecio")?.value;
    const disponibilidad = document.getElementById("filtroDisponibilidad")?.value;

    let obrasAMostrar = [];

    // le damos chance a que respire con await
    if (idGenero) {
        obrasAMostrar = await cargarObrasPorGenero(idGenero);
    } else if (disponibilidad === "disponible") {
        obrasAMostrar = await cargarObrasPorDisponibilidad();
    } else {
        obrasAMostrar = await obtenerCatalogoCompleto();
    }

    // Filtro secundario local para "No disponibles"
    if (disponibilidad === "no_disponible") {
        obrasAMostrar = obrasAMostrar.filter(o => o.status !== "DISPONIBLE");
    }

    // aplicamos el ordenamiento
    if (ordenPrecio === "asc") {
        obrasAMostrar.sort((a, b) => (a.precio || 0) - (b.precio || 0));
    } else if (ordenPrecio === "desc") {
        obrasAMostrar.sort((a, b) => (b.precio || 0) - (a.precio || 0));
    }

    // lo ordenamos localmente
    if ((idGenero || disponibilidad === "disponible") && ordenPrecio) {
        if (ordenPrecio === "asc") {
            obrasAMostrar.sort((a, b) => (a.precio || 0) - (b.precio || 0));
        } else if (ordenPrecio === "desc") {
            obrasAMostrar.sort((a, b) => (b.precio || 0) - (a.precio || 0));
        }
    }

    // ahora si tenemos la id del artista
    if (idArtista) {
        obrasAMostrar = obrasAMostrar.filter(o => o.artista && String(o.artista.id) === String(idArtista));
    }

    // llamado final para mostrar las obras
    mostrarObras(obrasAMostrar);
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
                <p>${ob.precio ? "$" + ob.precio : ""}</p>
            </div>
        `;
        gallery.appendChild(artCard);
    });
}