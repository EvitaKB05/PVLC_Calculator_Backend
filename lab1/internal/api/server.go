package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	log.Println("Starting server")

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("templates/*") // указываем путь, по которому лежат html/шаблоны, в моем случае это папка templates

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
