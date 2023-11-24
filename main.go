package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", pingpongHandler)

	r.Run(":8080")
}

func pingpongHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}
