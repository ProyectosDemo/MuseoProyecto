document.addEventListener('DOMContentLoaded', cargarObras);

async function cargarObras() {
    const query = `
        query {
            Obras(limit:6, offset:0) {
                id
                nombre
                foto
                artista {
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
        const container = document.getElementById("destacadas-container");
        if (!container) return;
        
        container.innerHTML = "";

        const obras = data.data?.Obras;
        if (!obras || obras.length === 0) {
            container.innerHTML = "<p>No hay obras destacadas disponibles</p>";
            return;
        }

        obras.forEach(ob => {
            const artCard = document.createElement("div");
            artCard.className = "art-card-destacada";

            // Imagen arriba a ancho completo, luego bloque de información
            artCard.innerHTML = `
                <a href="obra.html?id=${ob.id}">
                    <img src="${ob.foto}" alt="${ob.nombre}">
                </a>
                <div class="art-info">
                    <a href="obra.html?id=${ob.id}">
                        <h4>${ob.nombre}</h4>
                    </a>
                    <p>${ob.artista ? ob.artista.nombre : "Artista Desconocido"}</p>
                </div>
            `;
            container.appendChild(artCard);
        });

    } catch (err) {
        console.error("Error al cargar obras destacadas:", err);
        const container = document.getElementById("destacadas-container");
        if (container) {
            container.innerHTML = "<p>Error al cargar las obras destacadas</p>";
        }
    }
}