package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func getTodos(c *gin.Context, todos *[]Todo) {
	if *todos == nil {
		db := c.MustGet("db").(*sql.DB)
		getTodosFromDB(todos, db)
	}
	done := c.Query("done")
	var res []Todo
	switch done {
	case "true":
		res = Filter(
			todos,
			func(todo Todo) bool { return todo.Done },
		)
	case "false":
		res = Filter(
			todos,
			func(todo Todo) bool { return !todo.Done },
		)
	case "":
		res = *todos
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "done parameter should be either true or false"})
		return
	}
	c.JSON(http.StatusOK, res)
}
func getTodoByID(c *gin.Context, todos *[]Todo) {
	if *todos == nil {
		db := c.MustGet("db").(*sql.DB)
		getTodosFromDB(todos, db)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		notFound(c)
		return
	}
	for _, todo := range *todos {
		if todo.ID == id {
			c.JSON(http.StatusOK, todo)
			return
		}
	}
	notFound(c)
}

func postTodo(c *gin.Context, todos *[]Todo) {
	db := c.MustGet("db").(*sql.DB)
	if *todos == nil {
		getTodosFromDB(todos, db)
	}
	body, err := c.GetRawData()
	if err != nil {
		badRequest(c, "Request body is empty")
		return
	}
	var m map[string]any
	err = json.Unmarshal(body, &m)
	if err != nil {
		badRequest(c, "Body must contain a JSON object")
		return
	}

	title, ok := m["title"]
	if _, correctType := title.(string); !ok || !correctType {
		badRequest(c, "title parameter must be a string")
		return
	}
	var todo Todo
	todo.Title = title.(string)

	if done, ok := m["done"]; ok {
		if doneBool, ok := done.(bool); ok {
			todo.Done = doneBool
		} else {
			badRequest(c, "done parameter should be either true or false")
			return
		}
	}
	id := postTodoToDB(todo, db)
	todo.ID = id
	*todos = append(*todos, todo)
	c.JSON(http.StatusOK, todo)

}

func patchTodo(c *gin.Context, todos *[]Todo) {
	db := c.MustGet("db").(*sql.DB)
	if *todos == nil {
		getTodosFromDB(todos, db)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		notFound(c)
		return
	}
	var todo *Todo
	for index := range *todos {
		if (*todos)[index].ID == id {
			todo = &(*todos)[index]
		}
	}
	if todo == nil {
		notFound(c)
		return
	}
	body, err := c.GetRawData()
	if err != nil {
		badRequest(c, "Request body is empty")
		return
	}
	var m map[string]any
	err = json.Unmarshal(body, &m)
	if err != nil {
		badRequest(c, "Body must contain a JSON object")
		return
	}
	err = updateFromMap(&todo.Title, m, "title")
	if err != nil {
		badRequest(c, "title parameter must be a string")
		return
	}
	err = updateFromMap(&todo.Done, m, "done")
	if err != nil {
		badRequest(c, "done parameter must be a boolean")
		return
	}
	updateTodoInDB(*todo, db)
	c.JSON(http.StatusOK, *todo)
	return

}

func deleteTodo(c *gin.Context, todos *[]Todo) {
	db := c.MustGet("db").(*sql.DB)
	if *todos == nil {
		getTodosFromDB(todos, db)
	}
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		notFound(c)
		return
	}
	for index, todo := range *todos {
		if intId == todo.ID {
			*todos = append((*todos)[:index], (*todos)[index+1:]...)
			deleteTodoInDB(todo.ID, db)
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
			return
		}
	}
	notFound(c)
}

func main() {
	var todos = new([]Todo)
	db := dbConnection()
	router := gin.Default()

	router.Use(cors.Default())
	router.Use(func(context *gin.Context) {
		context.Set("db", db)
		context.Next()
	})
	router.GET("/todos", func(context *gin.Context) {
		getTodos(context, todos)
	})
	router.POST("/todos", func(context *gin.Context) {
		postTodo(context, todos)
	})
	router.GET("/todos/:id", func(context *gin.Context) {
		getTodoByID(context, todos)
	})
	router.DELETE("/todos/:id", func(context *gin.Context) {
		deleteTodo(context, todos)
	})
	router.PATCH("/todos/:id", func(context *gin.Context) {
		patchTodo(context, todos)
	})

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}
