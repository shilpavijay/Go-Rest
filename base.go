package main 

import (
	"fmt"
	"database/sql"
	"bytes"
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
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
	//sql.Open doesn't give an error if the machine is not reachable. Hence it is necessary to Ping.
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}

	type Emp struct {
		id int
		first_name string
		last_name string
	}

	router := gin.Default()

	//API HANDLERS

	// GET: Single Employee's details:
	router.GET("/emp/:id", func(c *gin.Context) {
		var (
				emp Emp
				result gin.H
			)
		id := c.Param("id")
		row := db.QueryRow("select id,first_name,last_name from emp where id = ?;",id)
		err := row.Scan(&emp.id,&emp.first_name,&emp.last_name)

		if err != nil {
			//if no results found send Nill
				result = gin.H{
						"result": nil,
						"count": 0,
						}
		} else {
			result = gin.H{
					"result": gin.H{
						"id": emp.id,
						"first_name": emp.first_name,
						"last_name": emp.last_name,
						},
					"count": 1,
			}
		}
		c.JSON(http.StatusOK, result)
	})


	//GET all Employees data:
	router.GET("/allemps", func(c *gin.Context) {
		var buffer bytes.Buffer
		var (
				emp Emp
				allEmps []Emp
			)
		rows, err := db.Query("select id,first_name,last_name from emp;")
		if err != nil {
			fmt.Print(err.Error())
		}

		for rows.Next() {
			err = rows.Scan(&emp.id, &emp.first_name,&emp.last_name)
			if err != nil {
				fmt.Print(err.Error())
			}
			allEmps = append(allEmps, emp)
			d,err := json.Marshal(emp)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Sprintf(string(d))

		}

		for i:=0;i<len(allEmps);i++ {
			buffer.WriteString("id: ")
			buffer.WriteString(string(allEmps[i].id))
			buffer.WriteString(" First Name: ")
			buffer.WriteString(allEmps[i].first_name)
			buffer.WriteString(" Last Name: ")
			buffer.WriteString(allEmps[i].last_name)
			buffer.WriteString("\n")
		}
		result := buffer.String()
		fmt.Println(result)
		defer rows.Close()
		c.JSON(http.StatusOK, gin.H{
				"result": result,
				"count": len(allEmps),
			})
		})


	//POST a new Employee data:
	router.POST("/emp", func(c *gin.Context) {
		var buffer bytes.Buffer
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("insert into emp(first_name,last_name) values (?,?);")
		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = stmt.Exec(first_name,last_name)
		if err != nil {
			fmt.Println(err.Error())
		}

		//append String to buffer - Faster, Performance boosted!
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf(" %s successfully inserted", name),
			})
	})


	//Delete resources:
	router.DELETE("/emp/:id", func(c *gin.Context) {
		id := c.Param("id")
		stmt, err := db.Prepare("delete from emp where id= ?;")
		if err != nil {
			fmt.Println(err.Error())
		}

		_,err = stmt.Exec(id)
		if err != nil {
			fmt.Println(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted userid: %s", id),
			})
	})


	//PUT - update an employee's data:
	router.PUT("/emp/:id", func(c *gin.Context) {
		var buffer bytes.Buffer
		id := c.Param("id")
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		stmt, err := db.Prepare("update emp set first_name= ?, last_name= ? where id= ?")
		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = stmt.Exec(first_name,last_name,id)
		if err != nil {
			fmt.Println(err.Error())
		}

		//append strings
		buffer.WriteString(first_name)
		buffer.WriteString(" ")
		buffer.WriteString(last_name)
		defer stmt.Close()
		name := buffer.String()
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully updated user %s", name),
			})
	})

	router.Run(":3000")
}

