package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// type Person struct {
// 	Name string `json:"name"`
// 	Age  uint   `json:"age"`
// }

type Person struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

const (
	hostname  = "127.0.0.1"
	host_port = 5432
	username  = "postgres"
	password  = ""
	database  = "postgres"
)

func opendatabaseConnectio() *sql.DB {
	pg_con_string := fmt.Sprintf("port=%d host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host_port, hostname, username, password, database)

	db, err := sql.Open("postgres", pg_con_string)
	if err != nil {
		panic(err)
		fmt.Print("NO CONNECTION")
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/people", AllPeopleHandler)
	router.HandleFunc("/newperson", AddPersonHandler)
	log.Fatal(http.ListenAndServe(":8080", router))
}

// AllPersonHandler retrives all person register
func AllPeopleHandler(w http.ResponseWriter, r *http.Request) {
	db := opendatabaseConnectio()
	var people []Person
	rows, err := db.Query("SELECT * FROM PERSON")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var person Person
		rows.Scan(&person.Name, &person.Nickname)
		people = append(people, person)
	}
	peopleBytes, _ := json.MarshalIndent(people, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
	defer db.Close()
}

func AddPersonHandler(w http.ResponseWriter, r *http.Request) {
	var person Person
	db := opendatabaseConnectio()
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	insertPersonSQL := `INSERT INTO person (name, nickname) VALUES ($1, $2)`
	_, err = db.Exec(insertPersonSQL, person.Name, person.Nickname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	db.Close()
}
