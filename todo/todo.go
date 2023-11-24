package todo

import (
	"net/http"
	"strings"

	auth "github.com/benzt2h/gin-todo-api/auh"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	Title string `json:"text" binding:"required"`
	gorm.Model
}

func (Todo) TableName() string {
	return "todolist"
}

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (th *TodoHandler) NewTask(ctx *gin.Context) {
	s := ctx.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(s, "Bearer ")

	if err := auth.Protect(tokenString); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var todo Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r := th.db.Create(&todo)
	if err := r.Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}
