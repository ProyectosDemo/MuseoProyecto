// 1. CREAR EL PANEL VISUAL (HUD) AUTOMÁTICAMENTE EN LA PANTALLA
document.addEventListener("DOMContentLoaded", () => {
    // Crear el contenedor principal del HUD
    const hud = document.createElement("div");
    hud.id = "db-fps-hud";
    hud.innerHTML = `
        <div class="hud-title">DB LATENCY (PING)</div>
        <div class="hud-item">MySQL: <span id="hud-mysql" class="hud-value">--</span></div>
        <div class="hud-item">MongoDB: <span id="hud-mongodb" class="hud-value">--</span></div>
        <div class="hud-item">Cassandra: <span id="hud-cassandra" class="hud-value">--</span></div>
        <div class="hud-item">Neo4j: <span id="hud-neo4j" class="hud-value">--</span></div>
    `;

    // Crear el botón flotante para alternar visualización
    const toggleBtn = document.createElement("button");
    toggleBtn.id = "db-fps-toggle";
    toggleBtn.innerText = "📊 DB FPS";

    // Aplicar los estilos integrados en la esquina superior izquierda
    const style = document.createElement("style");
    style.innerHTML = `
        #db-fps-hud {
            position: fixed; 
            top: 55px; 
            left: 15px; 
            background: rgba(10, 10, 10, 0.9); 
            color: #00ff66;
            font-family: 'Consolas', 'Courier New', monospace; 
            font-size: 11px;
            padding: 12px; 
            border-radius: 6px; 
            border: 1px solid #00ff66;
            z-index: 99999; 
            box-shadow: 0 0 15px rgba(0, 255, 102, 0.3); 
            min-width: 150px;
            transition: opacity 0.3s ease, transform 0.3s ease;
        }
        #db-fps-hud.hidden {
            opacity: 0;
            transform: translateY(-20px);
            pointer-events: none; 
        }
        #db-fps-toggle {
            position: fixed;
            top: 15px;
            left: 15px; 
            background: #0a0a0a;
            color: #00ff66;
            border: 1px solid #00ff66;
            font-family: 'Consolas', 'Courier New', monospace;
            font-size: 11px;
            padding: 6px 12px;
            border-radius: 4px;
            cursor: pointer;
            z-index: 100000; 
            box-shadow: 0 0 10px rgba(0, 255, 102, 0.2);
            transition: all 0.2s ease;
        }
        #db-fps-toggle:hover {
            background: #00ff66;
            color: #0a0a0a;
            box-shadow: 0 0 15px rgba(0, 255, 102, 0.5);
        }
        .hud-title { font-weight: bold; border-bottom: 1px solid #00ff66; margin-bottom: 6px; padding-bottom: 3px; color: #ffffff; }
        .hud-item { display: flex; justify-content: space-between; margin: 4px 0; }
        .hud-value { font-weight: bold; color: #00ffff; }
    `;

    // Adjuntar los elementos generados a la página
    document.head.appendChild(style);
    document.body.appendChild(toggleBtn);
    document.body.appendChild(hud);

    // Lógica del botón para ocultar/mostrar
    toggleBtn.addEventListener("click", () => {
        hud.classList.toggle("hidden");
    });
});

// 2. EL TRUCO PARA INTERCEPTAR TODOS LOS FETCH DE CUALQUIER ARCHIVO JS
(function() {
    const originalFetch = window.fetch;
    window.fetch = async function (...args) {
        const response = await originalFetch(...args);
        const clonedResponse = response.clone();

        if (clonedResponse.headers.get("content-type")?.includes("application/json")) {
            clonedResponse.json()
                .then(responseObj => {
                    if (responseObj && responseObj.extensions) {
                        const ext = responseObj.extensions;
                        
                        // Actualizar los textos de forma segura evaluando si existen
                        if (ext.db_latency_mysql) {
                            document.getElementById('hud-mysql').innerText = ext.db_latency_mysql;
                        }
                        if (ext.db_latency_mongodb) {
                            document.getElementById('hud-mongodb').innerText = ext.db_latency_mongodb;
                        }
                        if (ext.db_latency_cassandra) {
                            document.getElementById('hud-cassandra').innerText = ext.db_latency_cassandra;
                        }
                        if (ext.db_latency_neo4j) {
                            document.getElementById('hud-neo4j').innerText = ext.db_latency_neo4j;
                        }
                    }
                })
                .catch(() => {});
        }
        return response;
    };
})();