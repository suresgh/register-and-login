package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	UserID      int    `json:"userID"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	City        string `json:"city"`
	State       string `json:"state"`
	Password    string `json:"password"`
}

func main() {
	// Connect to the database
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
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

	router.DELETE("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		res, err := db.Exec("DELETE FROM usersinfo WHERE userid=$1", idInt)

		if err != nil {
			log.Fatal(err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}

		if rowsAffected > 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("User with ID %s deleted successfully", id)})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("User with ID %s not found", id)})
		}
	})

	router.PUT("/api/users/:id", func(c *gin.Context) {

		userID := c.Param("id")

		var updateData struct {
			Field string `json:"field"`
			Value string `json:"value"`
		}
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var sql string
		switch updateData.Field {
		case "username":
			sql = "UPDATE usersinfo SET username=$1 WHERE userId=$2"
		case "email":
			sql = "UPDATE usersinfo SET email=$1 WHERE userId=$2"
		case "phoneNumber":
			sql = "UPDATE usersinfo SET phoneNumber=$1 WHERE userId=$2"
		case "city":
			sql = "UPDATE usersinfo SET city=$1 WHERE userId=$2"
		case "state":
			sql = "UPDATE usersinfo SET state=$1 WHERE userId=$2"
		case "password":
			sql = "UPDATE usersinfo SET password=$1 WHERE userId=$2"
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid field"})
			return
		}
		_, err := db.Exec(sql, updateData.Value, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Field updated successfully"})
	})

	router.GET("/api/users", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM usersinfo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var users []Person
		for rows.Next() {
			var user Person
			if err := rows.Scan(&user.UserID, &user.Username, &user.Email, &user.PhoneNumber, &user.City, &user.State, &user.Password); err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, users)
	})

	router.Run(":8080")
}
