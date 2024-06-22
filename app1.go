package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
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

func main() {

	var json_file string = "users_data.json"
	var user_numbers int = 1000000
	var user_info [1000000]User

	file_input, er := os.Open(json_file)
	if er != nil {
		fmt.Println("Error during openning file")
	}

	fs, _ := file_input.Stat()
	var file_size int64 = fs.Size()

	buf_read_file := make([]byte, file_size+1)
	rn, _ := file_input.Read(buf_read_file)
	if rn <= 0 {
		fmt.Println("Error reading")
	}

	er2 := json.Unmarshal(buf_read_file, &user_info)
	if er2 != nil {
		fmt.Println("error dungin json unmarshaling")
	}

	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}

	stSql_userinfo, err := db.Prepare("INSERT INTO users(ID, Name, Email, PhoneNumber) VALUES(? , ? , ? , ?)")
	if err != nil {
		panic(err.Error())
	}

	stSql_addrinfo, err := db.Prepare("INSERT INTO users(ID, Street, City, State, Zipcode, Country) VALUES(?, ? , ? , ? , ? , ?)")
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < user_numbers; i++ {

		_, erSqli := stSql_userinfo.Exec("users", user_info[i].ID, user_info[i].Name, user_info[i].Email, user_info[i].PhoneNumber)
		if erSqli != nil {
			fmt.Println("Error during inserting users info of ", user_info[i].ID)
			return
		}

		num_address := len(user_info[i].Addresses)
		if num_address > 0 {
			for j := 0; j < num_address; j++ {

				_, erSqli_addr := stSql_addrinfo.Exec("addresses", user_info[i].ID, user_info[i].Addresses[j].Street,
					user_info[i].Addresses[j].City, user_info[i].Addresses[j].State, user_info[i].Addresses[j].ZipCode, user_info[i].Addresses[j].Country)
				if erSqli_addr != nil {
					fmt.Println("Error during inserting address of user ", user_info[i].ID)
					break
				}
			}

		}
	}

	db.Close()

}
