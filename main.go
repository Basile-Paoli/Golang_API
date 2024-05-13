package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TodoList struct {
	ID    int    `json:"id"`
	Title string `json:"title,omitempty"`
	Done  bool   `json:"done"`
}

var todos = []TodoList{
	{ID: 1, Title: "Courses", Done: false},
	{ID: 2, Title: "Sport", Done: false},
	{ID: 3, Title: "Vacances", Done: false},
}
var todosIncrement = 4

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}

func getTodos(c *gin.Context) {
	c.JSON(http.StatusOK, todos)
}
func getTodoByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		notFound(c)
		return
	}
	for _, todo := range todos {
		if todo.ID == id {
			c.JSON(http.StatusOK, todo)
			return
		}
	}
	notFound(c)
}

func postTodo(c *gin.Context) {
	var todo TodoList
	err := c.BindJSON(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	todo.Done = false
	todo.ID = todosIncrement
	todosIncrement += 1
	todos = append(todos, todo)
	c.JSON(http.StatusOK, todo)

}
func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		notFound(c)
		return
	}
	for index, item := range todos {
		if intId == item.ID {
			todos = append(todos[:index], todos[index+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
			return
		}
	}
	notFound(c)
}

func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodoByID)
	router.POST("/todos/", postTodo)
	router.DELETE("/todos/:id", deleteTodo)

	router.Run("localhost:8080")
}
