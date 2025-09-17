package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// connect to a database
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=testdb user=pradeek password=Deepakr_123")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	fmt.Println("Connected to database")

	// test my database
	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot pingged database")
	}

	fmt.Println("pinged from database")

	// get rows from table
	getAllRows(conn)

	// insert a row

	// query := `insert into users (first_name, last_name) values ($1, $2)`
	// _, err = conn.Exec(query, "Jack", "Brown")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// get rows from the table again
	// getAllRows(conn)

	// update a row
	// stmt := `update users set first_name = $1 where first_name = $2`
	// _, err = conn.Exec(stmt, "Sean", "John")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// getAllRows(conn)

	// get row by id
	query := `select id, first_name, last_name from users where id = $1`
	var id int
	var first_name, last_name string
	row := conn.QueryRow(query, 2)

	err = row.Scan(&id, &first_name, &last_name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id, first_name, last_name)

	// delete a row from the table
	// query = `delete from users where id = $1`
	// _, err = conn.Exec(query, 2)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// get row from the table
	// err = getAllRows(conn)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()
	var first_name, last_name string
	var id int
	for rows.Next() {
		err := rows.Scan(&id, &first_name, &last_name)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("Record is", id, first_name, last_name)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error scanning rows", err)
	}
	fmt.Println("------------------------------------")
	return nil
}
