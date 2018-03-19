# 2g_citas_ms

Microservicio encargado de la gestión de citas, para la aplicación dop.

Se usa como lenguaje de base Go, el cual nos facilita la creación de aplicaciones enfocadas en web y nos brinda herramientas para la creación de servicios web.

Como base de datos se usa mongoDb, se escogió debido a que cada cita puede ser variable, tanto en el número de personas como en los lugares donde se realizara, de esta forma podremos integrar un arreglo para estos campos, y consultarlos de forma mas sencilla.


Dependencias:

mux para facilitar el CRUD REST:
  go get github.com/gorilla/mux

mgo para la conexión a la BD mongo:
  go get gopkg.in/mgo.v2
