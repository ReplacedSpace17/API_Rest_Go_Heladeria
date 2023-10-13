package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Configura la conexión a la base de datos MySQL.
	var err error
	db, err = sql.Open("mysql", "root:Javier1234567890@tcp(heladeriago_mysql_1:3306)/Heladeria")

	//db, err = sql.Open("mysql", "root:Javier1234567890$@tcp(localhost:3306)/Heladeria")
	//db, err = sql.Open("mysql", "root:Javier1234567890$@tcp(mysql-container:3306)/Heladeria")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Inicializa el enrutador Gorilla Mux.
	router := mux.NewRouter()

	// Define la ruta "/Hola" que responde con "Hola, mundo" y consulta la base de datos.
	router.HandleFunc("/Hola", func(w http.ResponseWriter, r *http.Request) {
		// Realiza una consulta de prueba a la base de datos.
		var result string
		err := db.QueryRow("SELECT 'Hola, base de datos'").Scan(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Responde con el resultado de la consulta y "Hola, mundo".
		fmt.Fprintf(w, "Resultado de la consulta: %s\n", result)
		fmt.Fprintln(w, "Hola, mundo")
	}).Methods("GET")

	// Define las rutas para manipular sabores.
	router.HandleFunc("/api/sabores", CreateSabor).Methods("POST")
	router.HandleFunc("/api/sabores", GetSabores).Methods("GET")
	router.HandleFunc("/api/sabores/{id}", GetSabor).Methods("GET")
	router.HandleFunc("/api/sabores/{id}", UpdateSabor).Methods("PUT")
	router.HandleFunc("/api/sabores/{id}", DeleteSabor).Methods("DELETE")

	// Inicia el servidor HTTP en el puerto 8080.
	fmt.Println("Servidor escuchando en el puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Sabor representa un sabor de helado.
type Sabor struct {
	ID     int     `json:"id"`
	Sabor  string  `json:"sabor"`
	Precio float64 `json:"precio"`
}

// CreateSabor crea un nuevo sabor en la base de datos.
func CreateSabor(w http.ResponseWriter, r *http.Request) {
	var sabor Sabor
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sabor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Inserta el nuevo sabor en la base de datos.
	_, err := db.Exec("INSERT INTO sabores(sabor, precio) VALUES(?, ?)", sabor.Sabor, sabor.Precio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responde con el sabor creado en formato JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sabor)
}

// GetSabores obtiene todos los sabores de la base de datos.
func GetSabores(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, sabor, precio FROM sabores")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sabores []Sabor
	for rows.Next() {
		var sabor Sabor
		if err := rows.Scan(&sabor.ID, &sabor.Sabor, &sabor.Precio); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sabores = append(sabores, sabor)
	}

	// Responde con la lista de sabores en formato JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sabores)
}

// GetSabor obtiene un sabor por su ID.
func GetSabor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sabor Sabor
	err = db.QueryRow("SELECT id, sabor, precio FROM sabores WHERE id=?", id).Scan(&sabor.ID, &sabor.Sabor, &sabor.Precio)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Sabor no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Responde con el sabor en formato JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sabor)
}

// UpdateSabor actualiza un sabor por su ID.
func UpdateSabor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var sabor Sabor
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&sabor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Actualiza el sabor en la base de datos.
	_, err = db.Exec("UPDATE sabores SET sabor=?, precio=? WHERE id=?", sabor.Sabor, sabor.Precio, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responde con un mensaje de éxito.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"mensaje": "Sabor con ID %d actualizado correctamente"}`, id)
}

// DeleteSabor elimina un sabor por su ID.
func DeleteSabor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Elimina el sabor de la base de datos.
	_, err = db.Exec("DELETE FROM sabores WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Responde con un mensaje de éxito.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"mensaje": "Sabor con ID %d eliminado correctamente"}`, id)
}
