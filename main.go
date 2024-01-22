package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// Model data untuk representasi objek di database
type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var db *sql.DB

func init() {
	var err error

	// Koneksi ke database PostgreSQL
	db, err = sql.Open("postgres", "user=postgres password=Since2024. host= dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Uji koneksi ke database
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")
}

func main() {
	// Inisialisasi router menggunakan Gin
	router := gin.Default()

	// Menentukan handler untuk route
	router.GET("/people", GetPeople)
	router.GET("/people/:id", GetPerson)
	router.POST("/people", CreatePerson)

	// Menjalankan server HTTP di port 8080
	log.Fatal(router.Run(":8080"))
}

// Handler untuk mendapatkan semua orang
func GetPeople(c *gin.Context) {
	var people []Person

	// Mengambil data dari database
	rows, err := db.Query("SELECT * FROM people")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	// Membaca hasil query dan menambahkannya ke slice people
	for rows.Next() {
		var person Person
		err := rows.Scan(&person.ID, &person.Name, &person.Age)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		people = append(people, person)
	}

	// Mengembalikan response dalam bentuk JSON
	c.JSON(http.StatusOK, people)
}

// Handler untuk mendapatkan satu orang berdasarkan ID
func GetPerson(c *gin.Context) {
	id := c.Param("id")
	var person Person

	// Mengambil data dari database berdasarkan ID
	row := db.QueryRow("SELECT * FROM people WHERE id=$1", id)
	err := row.Scan(&person.ID, &person.Name, &person.Age)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	// Mengembalikan response dalam bentuk JSON
	c.JSON(http.StatusOK, person)
}

// Handler untuk membuat orang baru
func CreatePerson(c *gin.Context) {
	var person Person

	// Membaca data JSON dari request body
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Menyimpan data ke database
	_, err := db.Exec("INSERT INTO people (name, age) VALUES ($1, $2)", person.Name, person.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mengembalikan response
	c.JSON(http.StatusCreated, gin.H{"message": "Person created successfully"})
}
