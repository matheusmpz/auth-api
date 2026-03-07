package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/matheusmpz/auth-api/db"
)

func main() {
    // Inicializa o banco de dados e o pool antes de receber qualquer requisição
    database := db.GetDB()
    if database == nil {
        log.Fatal("Falha ao inicializar o banco de dados")
    }

    // Inicializando as rotas com a porta default do gin
    router := gin.Default()
    // Inicializando o cors para evitar conflito de rotas
    router.Use(cors.Default())

    // rota básica para retorno de uma mensagem utilizando o router do gin
    router.GET("/ping", func(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, gin.H{"message": "Alô"})
    })

    // tudo que vier depois disso não é executado, por isso ele fica no final
    router.Run(":8080")
}