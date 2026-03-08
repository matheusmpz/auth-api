package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/matheusmpz/auth-api/controllers"
	"github.com/matheusmpz/auth-api/db"
	"github.com/matheusmpz/auth-api/middlewares"
)

func main() {
    // Inicializa o banco de dados e o pool antes de receber qualquer requisição
    database := db.GetDB()
    if database == nil {
        log.Fatal("Falha ao inicializar o banco de dados")
    }

    // Inicializa controller
	userCtrl := &controllers.UserController{DB: database}

    // Inicializando as rotas com a porta default do gin
    router := gin.Default()
    // Inicializando o cors para evitar conflito de rotas
    router.Use(cors.Default())

    // iniciando as rotas da aplicação
    router.POST("/register", userCtrl.Register)
	router.POST("/login", userCtrl.Login)
	router.POST("/activate", userCtrl.Activate)

    auth := router.Group("/")
    auth.Use(middlewares.AuthMiddleware())
    {
        auth.GET("/users/:id", userCtrl.GetUser)
        auth.PUT("/users/:id", userCtrl.UpdateUser)
        auth.DELETE("/users/:id", userCtrl.DeleteUser)
        auth.PATCH("/users/:id/activate", userCtrl.ActivateUser)
        auth.PATCH("/users/:id/block", userCtrl.BlockUser)
        auth.PATCH("/users/:id/unblock", userCtrl.UnblockUser)
    }

    // tudo que vier depois disso não é executado, por isso ele fica no final
    router.Run(":8080")
}