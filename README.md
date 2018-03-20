# 2g_citas_ms dop

Microservicio encargado de la gestión de citas, para la aplicación dop.

Se usa como lenguaje de base Go, el cual nos facilita la creación de aplicaciones enfocadas en web y nos brinda herramientas para la creación de servicios web.

Como base de datos se usa mongodb, se escogió debido a que cada cita puede ser variable, tanto en el número de personas como en los lugares donde se realizara (itinerario), de esta forma podremos integrar un arreglo para estos campos, y consultarlos de forma mas sencilla.


### Dependencias:

1. mux para facilitar el CRUD REST:

  `go get github.com/gorilla/mux`

2. mgo para facilitar la conexión a mongodb:

  `go get gopkg.in/mgo.v2`


### Ejecución

En la carpeta del archivo ejecutar las siguientes instucciones:

1. Subir la bd a rancher:

`docker-compose up`

2. Subir el programa a Rancher

`docker build -t citas-ms .`

3. Iniciar el programa

`docker run --name citas-ms -p 3023:3300 citas-ms`

Con esto a tendremos en Rancher el servicio en ejecución

![Microservicio en Rancher](Imagenes/rancher.png "Microservicio en Rancher")

Ahora procedemos a ingresar a la dirección:
<http://192.168.99.101:3023> desde donde probaremos los diferentes servicios con la ayuda de postman.

CRUD | RUTA
:--: | :--
POST | /citas
GET  | /citas/{id}
GET  | /citas/personas/{id}
PUT  | /citas/{id}
DELETE | /citas/{id}
