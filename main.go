package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	auth "github.com/benzt2h/gin-todo-api/auh"
	"github.com/benzt2h/gin-todo-api/todo"
)

var (
	buildcommit = "dev"
	buildtime   = time.Now().String()
)

func main() {
	_, errLiveFile := os.Create("/tmp/live")
	if errLiveFile != nil {
		log.Fatal(errLiveFile)
	}
	defer os.Remove("/tmp/live")

	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("Please consider env varialble: %s", err)
	}

	db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()
	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	r.GET("/limitz", limitedHandler)
	r.GET("/x", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"buildcommit": buildcommit,
			"buildtime":   buildtime,
		})
	})

	r.GET("/tokenz", auth.AccessToken(os.Getenv("SIGN")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	todoHandler := todo.NewTodoHandler(db)
	protected.POST("/todos", todoHandler.NewTask)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shtting down gracfully, press Ctrl+c again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}

}

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(ctx *gin.Context) {
	if !limiter.Allow() {
		ctx.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
