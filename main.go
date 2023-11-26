package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	auth "github.com/benzt2h/gin-todo-api/auh"
	"github.com/benzt2h/gin-todo-api/todo"
)

func main() {

	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("Please consider env varialble: %s", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	todoHandler := todo.NewTodoHandler(db)
	protected.POST("/todos", todoHandler.NewTask)

	r.Run()
}
