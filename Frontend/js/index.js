// codigo de obras de arte


// Llamar a la función para cargar las obras al cargar la página
document.addEventListener('DOMContentLoaded', cargarObras);

async function cargarObras() {
    const query = `
        query {
            Obras(limit:4, offset:0) {
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
        const gallery = document.getElementById("gallery");
        gallery.innerHTML = "";

        const obras = data.data.Obras;
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
    </div>
`;
            gallery.appendChild(artCard);
        });

    } catch (err) {
        console.error(err);
        document.getElementById("gallery").innerHTML = "<p>Error cargando obras</p>";
    }
}
