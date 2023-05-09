package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "12345"
	dbname   = "mydatabase"
)

type Users struct {
	Userid      string `json:"userid"`
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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()

	router.GET("/excelusers", func(c *gin.Context) {
		var users []Users

		rows, err := db.Query("SELECT * FROM usersinfo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		// Iterate over the rows and add them to the users slice
		for rows.Next() {
			var user Users
			err := rows.Scan(&user.Userid, &user.Username, &user.Email, &user.PhoneNumber, &user.City, &user.State, &user.Password)
			if err != nil {
				log.Fatal(err)
			}
			users = append(users, user)
		}

		// Create a new Excel file
		f := xlsx.NewFile()

		// Create a new sheet
		sheetName := "Sheet1"
		sheet, err := f.AddSheet(sheetName)
		if err != nil {
			log.Fatal(err)
		}

		// Set the header row
		header := sheet.AddRow()
		header.AddCell().Value = "Userid"
		header.AddCell().Value = "Username"
		header.AddCell().Value = "Email"
		header.AddCell().Value = "PhoneNumber"
		header.AddCell().Value = "City"
		header.AddCell().Value = "State"
		header.AddCell().Value = "Password"

		// Set the data rows
		for _, user := range users {
			row := sheet.AddRow()
			row.AddCell().Value = user.Userid
			row.AddCell().Value = user.Username
			row.AddCell().Value = user.Email
			row.AddCell().Value = user.PhoneNumber
			row.AddCell().Value = user.City
			row.AddCell().Value = user.State
			row.AddCell().Value = user.Password
		}

		// Save the Excel file
		err = f.Save("usersSheet.xlsx")
		if err != nil {
			log.Fatal(err)
		}

		c.Header("Content-Disposition", "attachment; filename=usersSheet.xlsx")
		c.File("usersSheet.xlsx")

	})
	router.Run(":8080")
}
