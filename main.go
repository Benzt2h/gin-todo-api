package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	auth "github.com/benzt2h/gin-todo-api/auh"
	"github.com/benzt2h/gin-todo-api/todo"
)

func main() {

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()

	handler := todo.NewTodoHandler(db)

	r.POST("/todos", handler.NewTask)

	r.GET("/tokenz", auth.AccessToken)

	r.Run(":8080")
}
