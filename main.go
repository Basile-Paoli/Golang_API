package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func getTodos(c *gin.Context, todos *[]Todo) {

	done := c.Query("done")
	var res []Todo
	switch done {
	case "true":
		res = Filter(
			*todos,
			func(todo Todo) bool { return todo.Done },
		)
	case "false":
		res = Filter(
			*todos,
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

	done, ok := m["done"]
	if doneBool, correctType := done.(bool); ok && correctType && doneBool == true {
		todo.Done = true
	} else {
		todo.Done = false
	}

	id := postTodoToDB(todo)
	todo.ID = id
	*todos = append(*todos, todo)
	c.JSON(http.StatusOK, todo)

}

func patchTodo(c *gin.Context, todos *[]Todo) {
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
	updateTodoInDB(*todo)
	c.JSON(http.StatusOK, *todo)
	return

}

func deleteTodo(c *gin.Context, todos *[]Todo) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		notFound(c)
		return
	}
	for index, todo := range *todos {
		if intId == todo.ID {
			*todos = append((*todos)[:index], (*todos)[index+1:]...)
			deleteTodoInDB(todo.ID)
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
			return
		}
	}
	notFound(c)
}

func main() {
	var todos = new([]Todo)
	getTodosFromDB(todos)
	router := gin.Default()
	router.GET("/todos", func(context *gin.Context) {
		getTodos(context, todos)
	})
	router.POST("/todos/", func(context *gin.Context) {
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

	router.Run("localhost:8080")
}
