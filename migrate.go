package main 

import (
	"fmt"
	"database/sql"
	// "bytes"
	// "net/http"

	// "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql","root:python098@(127.0.0.1:3306)/gotest")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Connection established successfully")
	}
	defer db.Close()

	//Checking for connection:
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	//Create Table:
	stmt, err := db.Prepare("CREATE TABLE emp (id int NOT NULL AUTO_INCREMENT, first_name varchar(40), last_name varchar(40), PRIMARY KEY(id));")
	if err != nil {
		fmt.Println(err.Error())
	} 

	_,err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("emp table migration successful")
	}
}

