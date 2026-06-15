document.addEventListener('DOMContentLoaded', manejarRecomendaciones);

async function manejarRecomendaciones() {
    // 1. CAPTURA AUTOMÁTICA: Extraemos el ID usando la clave exacta que definiste en tu login
    const idClienteLogueado = localStorage.getItem("clienteId");

    if (!idClienteLogueado) {
        document.getElementById("gallery").innerHTML = "<p>Por favor, inicia sesión para ver tus recomendaciones personalizadas.</p>";
        return;
    }

    // 2. Traer las 10 recomendaciones de ese usuario específico
    const recomendaciones = await obtenerRecomendaciones(idClienteLogueado);
    mostrarObras(recomendaciones);
}

async function obtenerRecomendaciones(idCliente) {
    // Definimos la Query estructurada con la variable idCliente limpia
    const operacionGraphQL = {
        query: `
            query GetRecomendaciones($id: ID!) {
                obtenerRecomendaciones(idCliente: $id) {
                    id
                    nombre
                    foto
                }
            }
        `,
        variables: {
            id: String(idCliente)
        }
    };
    
    try {
        const res = await fetch("http://localhost:8080/query", { 
            method: "POST", 
            headers: { "Content-Type": "application/json" }, 
            body: JSON.stringify(operacionGraphQL) 
        });
        
        const data = await res.json();
        
        if (data.errors) {
            console.error("Errores de GraphQL:", data.errors);
            document.getElementById("gallery").innerHTML = "<p>Error procesando las recomendaciones</p>";
            return [];
        }

        return data.data.obtenerRecomendaciones || [];
    } catch (err) {
        console.error(err);
        document.getElementById("gallery").innerHTML = "<p>Error cargando obras recomendadas</p>";
        return [];
    }
}

// Renderiza usando tus clases nativas (.art-card, .art-info) para tu hoja de estilos
function mostrarObras(obras) {
    const gallery = document.getElementById("gallery");
    gallery.innerHTML = "";

    // Si el cliente no ha comprado nada, el backend de Neo4j devolverá una lista vacía
    if (!obras || obras.length === 0) {
        gallery.innerHTML = "<p>No tienes recomendaciones disponibles. ¡Debes comprar al menos una obra para que podamos conocer tus gustos!</p>";
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