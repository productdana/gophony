package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "dana_db"
)

func sanitizePhone(phone string) int {
	reg, err := regexp.Compile("[^0-9]+")

	if err != nil {
		log.Fatal("sanitizePhone err:", err)
	}

	sanitizedStr := reg.ReplaceAllString(phone, "")
	sanitizedPhone, _ := strconv.Atoi(sanitizedStr)

	return sanitizedPhone
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"dbname=%s sslmode=disable", host, port, user, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping() // open connection to db
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	// sqlStatement := `
	// INSERT INTO phone_numbers (phone_dirty)
	// VALUES ($1)`

	// dirtyNumbers := []string{
	// 	"123 456 7891",
	// 	"(123) 456 7892",
	// 	"(123) 456-7893",
	// 	"123-456-7894",
	// 	"123-456-7890",
	// 	"1234567892",
	// 	"(123)456-7892",
	// }

	// func addDirtyPhonesToTable()  {
	// 	for _, phone := range dirtyNumbers {
	// 		// add dirty phone to table
	// 		_, err := db.Exec(sqlStatement, phone)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }

	// for each record, sanitize for numbers only and store in phone column
	rows, err := db.Query("SELECT * FROM phone_numbers LIMIT $1", 10)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          int
			phone       int
			phone_dirty string
		)
		err := rows.Scan(&id, &phone, &phone_dirty)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("id:", id, "phonedirty:", phone_dirty)
		sanitizedPhone := sanitizePhone(phone_dirty)
		sqlStatement := `
		UPDATE phone_numbers
		SET phone=$2
		WHERE id=$1;
		`
		_, err = db.Exec(sqlStatement, id, sanitizedPhone)
		if err != nil {
			panic(err)
		}
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
	rows.Close()
}
