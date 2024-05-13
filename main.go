package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"log"
	"net/http"
	"strconv"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title,omitempty"`
	Done  bool   `json:"done"`
}
type TodoList []Todo

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
}

func getTodos(c *gin.Context, todos *TodoList) {
	if *todos == nil {
		getTodosFromDB(todos)
	}

	done := c.Query("done")
	res := []Todo{}
	switch done {
	case "true":
		for _, todo := range *todos {
			if todo.Done {
				res = append(res, todo)
			}
		}
	case "false":
		for _, todo := range *todos {
			if !todo.Done {
				res = append(res, todo)
			}
		}
	case "":
		res = *todos
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "done parameter should be either true or false"})
		return
	}
	c.JSON(http.StatusOK, res)
}
func getTodoByID(c *gin.Context, todos *TodoList) {
	if *todos == nil {
		getTodosFromDB(todos)
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

func postTodo(c *gin.Context, todos *TodoList) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error reading body"})
		return
	}
	var m map[string]any
	err = json.Unmarshal(body, &m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	title, ok := m["title"]
	if _, correctType := title.(string); !ok || !correctType {
		c.JSON(http.StatusBadRequest, gin.H{"message": "title parameter must be a string"})
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
	if *todos == nil {
		getTodosFromDB(todos)
	} else {
		*todos = append(*todos, todo)
	}
	c.JSON(http.StatusOK, todo)

}

func patchTodo(c *gin.Context, todos *TodoList) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		notFound(c)
		return
	}
	for index := range *todos {
		if (*todos)[index].ID == id {
			err := c.BindJSON(&(*todos)[index])
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusOK, (*todos)[index])
			return
		}
	}
	notFound(c)
}

func deleteTodo(c *gin.Context, todos *TodoList) {
	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		notFound(c)
		return
	}
	for index, item := range *todos {
		if intId == item.ID {
			*todos = append((*todos)[:index], (*todos)[index+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
			return
		}
	}
	notFound(c)
}

func todoList(todos *TodoList) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("todos", todos)
		c.Next()
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var todos = new(TodoList)
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
