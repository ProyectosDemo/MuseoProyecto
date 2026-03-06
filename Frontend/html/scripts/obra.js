function listarObrasdDestacadas() {
    fetch('http://localhost:8080/obras/destacadas')
        .then(response => response.json())
        .then(data => {
            const obrasContainer = document.getElementById('obras-destacadas');
            obrasContainer.innerHTML = ''; // Limpiar el contenedor antes de agregar nuevas obras

            data.forEach(obra => {
                const obraElement = document.createElement('div');
                obraElement.classList.add('obra');

                obraElement.innerHTML = `
                    <h3>${obra.titulo}</h3>
                    <p>${obra.descripcion}</p>
                    <p><strong>Autor:</strong> ${obra.autor}</p>
                    <p><strong>Año:</strong> ${obra.anio}</p>
                `;

                obrasContainer.appendChild(obraElement);
            });
        })
        .catch(error => console.error('Error al cargar las obras destacadas:', error));
}

// Llamar a la función para cargar las obras destacadas al cargar la página
document.addEventListener('DOMContentLoaded', listarObrasdDestacadas);