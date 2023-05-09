package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "mydatabase"
)

type Person struct {
	Username    string `json:"Username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	City        string `json:"city"`
	State       string `json:"state"`
	Password    string `json:"password"`
}

func main() {

	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s ",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

	})

	router.POST("/people", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sql := "INSERT INTO usersinfo (username, email, phoneNumber, city, state, password) VALUES ($1, $2, $3, $4, $5, $6)"
		_, err := db.Exec(sql, person.Username, person.Email, person.PhoneNumber, person.City, person.State, person.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Person created successfully"})

	})

	router.Run(":8082")
}
