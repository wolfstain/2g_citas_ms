package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// hosts      = "dockercompose_mongodb_1:27017"
	hosts      = "localhost:27017"
	database   = "2g_citas_bd"
	username   = ""
	password   = ""
	collection = "citas"
)

// Cita crea el tipo de cita, para los JSON
type Cita struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Cita     string        `bson:"Cita" json:"Cita,omitempty"`
	Lugar    string        `bson:"Lugar" json:"Lugar,omitempty"`
	Fecha    string        `bson:"Fecha" json:"Fecha,omitempty"`
	Personas []string      `bson:"Personas" json:"Personas,omitempty"`
	Estado   string        `bson:"Estado" json:"Estado,omitempty"`
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

	cita.ID = bson.NewObjectId()

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

	col.FindId(bson.ObjectIdHex(params["id"])).One(&cita)
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

	col.Find(bson.M{"Personas": params["id"]}).All(&citas)
	jsonString, err := json.Marshal(citas)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(jsonString))
}

// EditCitaEndpoint Edita la cita dado su id y un json con la info req to w
func EditCitaEndpoint(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "EditCitaEndpoint: Edita la cita dado su id y un json con la info, no esta implementado, aun!")
	params := mux.Vars(req)
	col := mongoStore.session.DB(database).C(collection)
	var cita Cita

	json.NewDecoder(req.Body).Decode(&cita)
	err := col.UpdateId(bson.ObjectIdHex(params["id"]), &cita)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, "result : success")
}

// DeleteCitaEndpoint Elimina una cita dado su id req to w
func DeleteCitaEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	col := mongoStore.session.DB(database).C(collection)

	err := col.RemoveId(bson.ObjectIdHex(params["id"]))
	if err != nil {
		panic(err)
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

	log.Fatal(http.ListenAndServe("localhost:3300", router))
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
	return
}
