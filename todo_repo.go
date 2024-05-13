package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func dbConnection() *sql.DB {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	connStr := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPassword + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func getTodosFromDB(todos *TodoList) {
	db := dbConnection()
	defer db.Close()
	rows, err := db.Query("SELECT id,title,done FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	*todos = TodoList{}
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Done)
		if err != nil {
			log.Fatal(err)
		}
		*todos = append(*todos, todo)
	}
}

func postTodoToDB(todo Todo) int {
	db := dbConnection()
	defer db.Close()
	row := db.QueryRow("INSERT INTO todos(done, title) VALUES ($1,$2) RETURNING id", todo.Done, todo.Title)
	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}
