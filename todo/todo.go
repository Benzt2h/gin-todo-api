package todo

import (
	"log"
	"net/http"
	"strconv"

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
	var todo Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if todo.Title == "sleep" {
		transactionID := ctx.Request.Header.Get("TransactionID")
		aud, _ := ctx.Get("aud")
		log.Println(transactionID, aud, "not allowed")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "not allowed",
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

func (th *TodoHandler) List(ctx *gin.Context) {
	var todos []Todo
	r := th.db.Find(&todos)
	if err := r.Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, todos)

}

func (th *TodoHandler) Remove(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r := th.db.Delete(&Todo{}, id)
	if err := r.Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}
