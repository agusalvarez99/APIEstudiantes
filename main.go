package main

import (
	"APIEstudiantes/middleware"
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Estudiante struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	FechaNac string `json:"fecha_nac"`
	Correo   string `json:"correo"`
}

type estudiantes []Estudiante

var listaEstudiantes = estudiantes{
	{
		ID:       1,
		Nombre:   "Agus",
		Apellido: "Alvarez",
		FechaNac: "16/07/1999",
		Correo:   "agus@mail.com",
	},
	{
		ID:       2,
		Nombre:   "Martin",
		Apellido: "Toledo",
		FechaNac: "10/11/1999",
		Correo:   "martin@mail.com",
	},
}

func getEstudiantes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(listaEstudiantes)
}

func getEstudiante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	estudianteID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID de estudiante no valido!")
		return
	}

	for _, estudiante := range listaEstudiantes {
		if estudiante.ID == estudianteID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(estudiante)
		}
	}

}

func addEstudiante(w http.ResponseWriter, r *http.Request) {
	var newEstudiante Estudiante
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Datos invalidos!")
	}

	json.Unmarshal(reqBody, &newEstudiante)

	newEstudiante.ID = len(listaEstudiantes) + 1
	listaEstudiantes = append(listaEstudiantes, newEstudiante)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEstudiante)
}

func deleteEstudiante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	estudianteID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID de estudiante no valido!")
		return
	}

	for i, estudiante := range listaEstudiantes {
		if estudiante.ID == estudianteID {
			listaEstudiantes = append(listaEstudiantes[:i], listaEstudiantes[i+1:]...)
			fmt.Fprintf(w, "Estudiante eliminado con exito!")
		}
	}
}

func updateEstudiante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	estudianteID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "ID de estudiante no valido!")
		return
	}

	var updatedEstudiante Estudiante

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Datos no validos")
		return
	}

	json.Unmarshal(reqBody, &updatedEstudiante)

	for i, estudiante := range listaEstudiantes {
		if estudiante.ID == estudianteID {
			listaEstudiantes = append(listaEstudiantes[:i], listaEstudiantes[i+1:]...)
			updatedEstudiante.ID = estudianteID
			listaEstudiantes = append(listaEstudiantes, updatedEstudiante)
			fmt.Fprintf(w, "Estudiante actualizado correctamente!")
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Bienvenido a nuestra API!"}`))
	//fmt.Fprintf(w, "Bienvenido a nuestra API!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	//router.HandleFunc("/estudiantes", getEstudiantes).Methods("GET")
	router.Handle("/estudiantes", middleware.EnsureValidToken()(http.HandlerFunc(getEstudiantes)))

	//ruteo anterior
	router.HandleFunc("/estudiantes", addEstudiante).Methods("POST")
	router.HandleFunc("/estudiantes/{id}", getEstudiante).Methods("GET")
	router.HandleFunc("/estudiantes/{id}", deleteEstudiante).Methods("DELETE")
	router.HandleFunc("/estudiantes/{id}", updateEstudiante).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", router))
}
