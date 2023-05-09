package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "mydatabase"
)

type User struct {
	Userid      int    `json:"userid"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	City        string `json:"city"`
	State       string `json:"state"`
	Password    string `json:"password"`
}

func main() {
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
	})

	r.POST("/users", func(c *gin.Context) {
		var users []User
		file, err := c.FormFile("userdata")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		xlFile, err := xlsx.OpenFile(file.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows {

				user := User{
					Username:    row.Cells[0].String(),
					Email:       row.Cells[1].String(),
					PhoneNumber: row.Cells[2].String(),
					City:        row.Cells[3].String(),
					State:       row.Cells[4].String(),
					Password:    row.Cells[5].String(),
				}

				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM usersinfo WHERE email = $1", user.Email).Scan(&count)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				if count == 0 {

					_, err = db.Exec("INSERT INTO usersinfo (username, email, phonenumber, city, state, password) VALUES ($1, $2, $3, $4, $5, $6)", user.Username, user.Email, user.PhoneNumber, user.City, user.State, user.Password)
					if err != nil {
						if e, ok := err.(*pq.Error); ok {
							log.Println(e.Code.Name())
						}
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
				}
				users = append(users, user)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d users added successfully", len(users))})
	})

	r.DELETE("/users/delete", func(c *gin.Context) {
		file, err := c.FormFile("userdata")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		xlFile, err := xlsx.OpenFile(file.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var numDeleted int

		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows {

				user := User{
					Username:    row.Cells[0].String(),
					Email:       row.Cells[1].String(),
					PhoneNumber: row.Cells[2].String(),
					City:        row.Cells[3].String(),
					State:       row.Cells[4].String(),
					Password:    row.Cells[5].String(),
				}

				var count int
				err = db.QueryRow("SELECT COUNT(*) FROM usersinfo WHERE email = $1", user.Email).Scan(&count)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				if count > 0 {

					_, err = db.Exec("DELETE FROM usersinfo WHERE email = $1", user.Email)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					numDeleted++
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d users deleted successfully", numDeleted)})
	})

	r.Run(":8081")
}
