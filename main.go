package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model
	Name string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model
	Title string
	Author string
	CallNumber int `gorm:"unique_index"`
	PersonID int
}

var (
	person = &Person { Name: "Jack", Email: "jack@email.com", }
	books  = []Book{
		{Title: "The Rules of Thinking", Author: "Richard Templar", CallNumber: 1234, PersonID: 1},
		{Title: "Winnie The Pooh", Author: "Alan A. Miln", CallNumber: 2345, PersonID: 1},
	}
)

var db *gorm.DB
var err error

func main() {
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Failed to get environment variables")
	}
	dialect  := myEnv["DIALECT"]
	host     := myEnv["HOST"]
	dbPort   := myEnv["DBPORT"]
	user     := myEnv["USER"]
	dbName   := myEnv["NAME"]
	password := myEnv["PASSWORD"]

	// Loading environment variables
	// To use this shall run command "source .env" and import "os" package
	/*
	dialect  := os.Getenv("DIALECT")
	host     := os.Getenv("HOST")
	dbPort   := os.Getenv("DBPORT")
	user     := os.Getenv("USER")
	dbName   := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")
	*/

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	// Opening connection to database
	db, err = gorm.Open(dialect,dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	// Close connection to database when the main function finishes
	defer db.Close()

	// Make migrations to the database if they have not already been created
	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	/*
	db.Create(&person)
	for idx := range books {
		db.Create(&books[idx])
	}
	*/

	// API routes
	router := mux.NewRouter()
	router.HandleFunc("/people",getPeople).Methods("GET")
	router.HandleFunc("/person/{id}",getPerson).Methods("GET")
	router.HandleFunc("/create/person",createPerson).Methods("POST")
	router.HandleFunc("/delete/person/{id}",deletePerson).Methods("DELETE")
	router.HandleFunc("/books",getBooks).Methods("GET")
	router.HandleFunc("/book/{id}",getBook).Methods("GET")
	router.HandleFunc("/create/book",createBook).Methods("POST")
	router.HandleFunc("/delete/book/{id}",deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))

}

// API Controllers

// People Controllers

func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	db.Find(&people)
	json.NewEncoder(w).Encode(&people)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	var books []Book
	db.First(&person, params["id"])
	db.Model(&person).Related(&books)
	person.Books = books
	json.NewEncoder(w).Encode(&person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	createdPerson := db.Create(&person)
	err = createdPerson.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&person)
	}
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	db.First(&person, params["id"])
	db.Delete(&person)
	json.NewEncoder(w).Encode(&person)
}

// Book Controllers

func getBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book
	db.Find(&books)
	json.NewEncoder(w).Encode(&books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	db.First(&book, params["id"])
	json.NewEncoder(w).Encode(&book)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	createdBook := db.Create(&book)
	err = createdBook.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&book)
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book
	db.First(&book, params["id"])
	db.Delete(&book)
	json.NewEncoder(w).Encode(&book)
}

