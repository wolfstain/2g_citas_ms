package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	ai "github.com/night-codes/mgo-ai"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// hosts = "localhost:27017"
	hosts      = "192.168.99.101:27017"
	database   = "2g_citas_bd"
	username   = ""
	password   = ""
	collection = "citas"
)

// Cita crea el tipo de cita, para los JSON
type Cita struct {
	ID          int64  `bson:"_id,omitempty" json:"id,omitempty"`
	Cita        string `bson:"cita" json:"cita,omitempty"`
	Lugar       int    `bson:"lugar" json:"lugar,omitempty"`
	Fecha       string `bson:"fecha" json:"fecha,omitempty"`
	Personas    []int  `bson:"personas" json:"personas,omitempty"`
	Estado      string `bson:"estado" json:"estado,omitempty"`
	Visibilidad bool   `bson:"visibilidad" json:"visibilidad"`
}

// MongoStore crea el tipo de dato MongoStore
type MongoStore struct {
	session *mgo.Session
}

var mongoStore = MongoStore{}

// Endpoints

// CreateCitaEndpoint Crea la cita dado su id y un json con la info req to w el _d es autogenerado por mongo
func CreateCitaEndpoint(w http.ResponseWriter, req *http.Request) {
	col := mongoStore.session.DB(database).C(collection)

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		panic(err)
	}

	var cita Cita
	err = json.Unmarshal(b, &cita)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// cita.ID = bson.NewObjectId()
	cita.ID = int64(ai.Next(collection))

	err = col.Insert(&cita)
	if err != nil {
		panic(err)
	}

	jsonString, err := json.Marshal(&cita)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")

	w.Write(jsonString)
}

// GetCitaEndpoint writes the JSON encoding of req to w.
func GetCitaEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	col := mongoStore.session.DB(database).C(collection)
	var cita Cita

	idc, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprint(w, "ID es formato invalido")
		return
	}

	if col.FindId(idc).One(&cita) != nil {
		fmt.Fprint(w, "ID no encontrado")
		return
	}

	jsonString, err := json.Marshal(cita)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(jsonString))
}

// GetCitaPersonaEndpoint Devuelve la información de la cita dado el id de la persona req to w.
func GetCitaPersonaEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	col := mongoStore.session.DB(database).C(collection)
	citas := []Cita{}

	idp, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Printf("ID invalido")
		return
	}

	col.Find(bson.M{"personas": idp}).All(&citas)
	jsonString, err := json.Marshal(citas)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(jsonString))
}

// EditCitaEndpoint Edita la cita dado su id y un json con la info req to w
func EditCitaEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	col := mongoStore.session.DB(database).C(collection)
	var cita Cita

	idc, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprint(w, "ID es formato invalido")
		return
	}

	json.NewDecoder(req.Body).Decode(&cita)
	if col.UpdateId(idc, &cita) != nil {
		fmt.Fprint(w, "ID no encontrado")
		return
	}

	fmt.Fprint(w, "result : success")
}

// DeleteCitaEndpoint Elimina una cita dado su id req to w
func DeleteCitaEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	col := mongoStore.session.DB(database).C(collection)

	idc, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Fprint(w, "ID es formato invalido")
		return
	}

	if col.RemoveId(idc) != nil {
		fmt.Fprint(w, "ID no encontrado")
		return
	}

	fmt.Fprint(w, "result : success")
}

func main() {
	//Crea la sesion de MongoDB
	session := initialiseMongo()
	mongoStore.session = session

	router := mux.NewRouter()

	// endpoints "github.com/gorilla/mux"
	router.HandleFunc("/citas", CreateCitaEndpoint).Methods("POST")
	router.HandleFunc("/citas/{id}", GetCitaEndpoint).Methods("GET")
	router.HandleFunc("/citas/personas/{id}", GetCitaPersonaEndpoint).Methods("GET")
	router.HandleFunc("/citas/{id}", EditCitaEndpoint).Methods("PUT")
	router.HandleFunc("/citas/{id}", DeleteCitaEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3300", router))
}

// initialiseMongo, se encarga de inicializar mongo con los datos almacencados como constantes
func initialiseMongo() (session *mgo.Session) {
	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	// connect AutoIncrement to collection "counters"
	ai.Connect(session.DB(database).C("counters"))
	return
}
