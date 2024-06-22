package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Addresses   []Address `json:"addresses"`
}

func getInfo(w http.ResponseWriter, req *http.Request) {

	var user_x User
	var addr_x Address

	user_id := strings.TrimPrefix(req.URL.Path, "/api/")

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("SELECT ID, Name, Email, PhoneNumber FROM users WHERE ID = ?")
	if err != nil {
		log.Fatal(err)
	}

	err2 := stmt.QueryRow(user_id).Scan(&user_x.ID, &user_x.Name, &user_x.Email, &user_x.PhoneNumber)
	if err2 != nil {
		log.Fatal(err)
	}

	//show it
	fmt.Fprintf(w, "User info>> ID: %s , Name: %s , Email: %s , PhoneNumber: %s", user_x.ID, user_x.Name, user_x.Email, user_x.PhoneNumber)

	rows, err := db.Query("SELECT Street, City, State , ZipCode, Country FROM addresses WHERE ID = ?", user_id)
	if err2 != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&addr_x.Street, &addr_x.City, &addr_x.State, &addr_x.ZipCode, &addr_x.Country)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Fprintf(w, "User Address>> Street: %s , City: %s , State: %s , ZipCode: %s, Country: %s", addr_x.Street, addr_x.City, addr_x.State, addr_x.ZipCode, addr_x.Country)
	}
}

func main() {

	http.HandleFunc("/api/", getInfo)

	http.ListenAndServe(":80", nil)

}
