package main

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func dbConnection() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	connStr := "user=" + dbUser + " dbname=" + dbName + " password=" + dbPassword + " sslmode=disable" + " port=" + dbPort + " host=" + dbHost
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func getTodosFromDB(todos *[]Todo, db *sql.DB) {
	rows, err := db.Query("SELECT id,title,done FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	*todos = []Todo{}
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Done)
		if err != nil {
			log.Fatal(err)
		}
		*todos = append(*todos, todo)
	}
}

func postTodoToDB(todo Todo, db *sql.DB) int {
	row := db.QueryRow("INSERT INTO todos(done, title) VALUES ($1,$2) RETURNING id", todo.Done, todo.Title)
	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}
func updateTodoInDB(todo Todo, db *sql.DB) Todo {
	row := db.QueryRow("UPDATE todos set done = $1, title=$2 where id = $3 returning id,title,done", todo.Done, todo.Title, todo.ID)
	var res Todo
	err := row.Scan(&res.ID, &res.Title, &res.Done)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
func deleteTodoInDB(id int, db *sql.DB) {
	_, err := db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
}
