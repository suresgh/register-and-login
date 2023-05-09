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

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s ", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

	})

	router.GET("/api", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		var count int

		err := db.QueryRow("SELECT COUNT(*) FROM usersinfo WHERE username= $1 AND password= $2", username, password).Scan(&count)

		if err != nil {
			log.Fatal(err)
		}

		if count > 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "login successfully"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Invalid username or password"})
		}
	})

	router.Run(":8081")
}
