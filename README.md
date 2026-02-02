El Museo de Arte Contemporáneo le ha encomendado que realice un sitio web que
permita la venta de las obras en exposición.  

Cada obra expuesta tiene un nombre y está asociada a un género (por ahora: pintura,
escultura, fotografía, cerámica y orfebrería) y a su artista creador. Además, tiene
un precio de venta y una fecha de creación. Es de interés conocer otros datos de la
obra  como por ejemplo en caso de las esculturas: material, peso, dimensiones (largo,
ancho, profundidad), etc. Se debe indagar las características de interés de acuerdo
al tipo  de género.  Cada artista puede dedicarse a más de un género.  

Se debe mostrar una foto de la obra expuesta y el nombre del autor de la obra debe
contener un link que lleve a los datos del autor (se debe mostrar una pequeña
biografia o datos de interés: fecha de nacimiento, nacionalidad, fotografía, entre
otros)  

Se debe permitir al visitante del sitio buscar obras de acuerdo a los siguientes
criterios:  
• De acuerdo a un género seleccionado por el  
• De acuerdo a un artista  
• Obras disponibles en venta de acuerdo al precio (de menor a mayor)  

Para que un visitante del sitio pueda adquirir una obra primero debe registrarse
como comprador. Para ello debe llenar un formulario con sus datos personales y de
una tarjera de crédito a la cual se le cargara el monto de 10 dolares por membresía. 
En retorno el sistema generara una clave aleatoria (código de seguridad) que será 
enviada al correo electrónico suministrado. Cada comprador deberá tener su login y 
password

En caso de extravío del código de seguridad, el usuario deberá responder a tres
preguntas de seguridad para que le pueda ser enviado otro código al correo. Al
momento de presionar el botón de compra el sistema deberá solicitar al comprador
su código de seguridad. La obra pasará del estatus Disponible a Reservada.  

Al emitirse la factura se considera venta concretada.  Se deben almacenar los datos
para el envio de la obra (el museo se encarga del envio)  

La venta de cada obra se hace de forma individual, es decir, cada factura debe incluir  
solo una obra y debe tener asociado el nombre de la obra vendida, el precio, datos
de pago y fecha de la venta. Esta factura es hecha por cualquier administrador del
sistema, una vez concretada la venta con el comprador. Los precios de las obras NO
incluyen el IVA, por lo tanto debe ser calculado en la venta. Por cada venta el museo
se queda con un porcentaje de ganancia previamente establecido con el artista, que
puede ir entre el 5 y 10%.  

Existen varios empleados en el museo que pueden tener el estatus de usuario
administrador. Cada uno de ellos deberá tener un login y password para poder entrar  
al sistema. 
Dentro del sistema cualquier usuario administrador podrá realizar las  
siguientes actividades:  
- CRUD de la data  
- Facturacion  
- Consultas:  
o Listado de obras vendidas en un periodo dado por el usuario  
o Resumen de facturación dado un periodo (código de factura, fecha, precio de la
obra, ganancia del museo (en porcentaje y en dólares) , total recaudado) 
o Resumen de membresías dado un período  


 
Un trabajador del museo se comunicará con el posible comprador en las próximas 24  
horas para concretar los detalles de la venta. De concretarse la venta, la obra pasará  
al estatus de Vendida, de lo contrario volverá a su estatus inicial (Disponible). 
