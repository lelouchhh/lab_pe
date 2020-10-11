package main

import (
	"database/sql"
	"fmt"
	_"fmt"
	"log"
	_"log"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func main() {
	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil{
		log.Println("error in db")
	}
	var ofirstName, oLastName string
	fmt.Println("enter your first name: ")
	fmt.Scan(&ofirstName)
	fmt.Println("enter your last name: ")
	fmt.Scan(&oLastName)
	if (!dataExisting(sqliteDatabase, ofirstName, oLastName)){
		fmt.Println("we add your data in data base")
		CreateUser(sqliteDatabase, ofirstName, oLastName)
	}
	defer sqliteDatabase.Close()
	fmt.Println("wait a bit..\nwe check your subscribe!")
	s := subscribeVerification(sqliteDatabase, ofirstName, oLastName)
	if !s{
		log.Println("your enter limit is end")
		return
	}
	fmt.Printf("Hello mister %s %s\n", ofirstName, oLastName)
	_, lim := EditUseLimit(sqliteDatabase, ofirstName, oLastName, s)
	fmt.Printf("your limit is: %d\n", lim)
}

func subscribeVerification(db *sql.DB, firstName, lastName string) bool{

	var col int
	sqlStmt := `SELECT use_limit FROM users WHERE first_name = ? AND last_name = ?`
	row := db.QueryRow(sqlStmt, firstName, lastName).Scan(&col)
	if row == sql.ErrNoRows{
		log.Println(row)
	}
	if col <= 0 {
		return false
	}
	return true
}

func dataExisting(db *sql.DB, firstName,lastName string) bool{
	var id int
	sqlStmt := `SELECT ID FROM users WHERE first_name = ? AND last_name = ?`
	err := db.QueryRow(sqlStmt, firstName, lastName).Scan(&id)
	if err == sql.ErrNoRows{
		log.Println(err)
		log.Printf("%s %s is not existing in DATA BASE", lastName, firstName)
		return false
	}
	return true
}
func CreateUser(db *sql.DB, firstName, lastName string) bool{
	sqlStatement := `
		INSERT INTO users (first_name, last_name, use_limit)
		VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, firstName, lastName, 4)
	if err != nil {
		log.Println(err)
	  	log.Printf("%s %s is not created", firstName, lastName)
		return false
	}

	return true
}

func EditUseLimit(db *sql.DB, firstName, lastName string, stat bool) (bool, int){
	if !stat{
		return false, 0
	}
	var col int
	sqlStmt := `SELECT use_limit FROM users WHERE first_name = ? AND last_name = ?`
	row := db.QueryRow(sqlStmt, firstName, lastName).Scan(&col)
	if row == sql.ErrNoRows{
		log.Println(row)
		log.Printf("%s %s can't return use_limit", lastName, firstName)
	}
	col--
	sqlStatement := `
		UPDATE users
		SET use_limit = $1 
		WHERE first_name = $2 and last_name = $3;`
	_, err := db.Exec(sqlStatement, col, firstName, lastName)
	if err != nil {
		log.Println(err)
		return false, 0
	}
	return true, col
}

