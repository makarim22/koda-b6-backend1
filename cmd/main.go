package main

import (
    "golang-backend/internal/handler"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "golang-backend/docs"
)


// @title           Simple API
// @version         1.0
// @description     Simple API to learn Swaggo
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8888
// @BasePath  /
func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.Data(200, "text/plain", []byte("Hello!"))
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// User Endpoints
	r.GET("/users", handler.GetUsers)
	r.GET("/users/:id", handler.GetUserByID)
	r.POST("/users", handler.CreateUser)
	r.PUT("/users/:id", handler.UpdateUser)

	// Auth Endpoints
	r.POST("/auth/register", handler.RegisterUser)
	r.POST("/auth/login", handler.LoginUser)

	// Product Endpoints
	r.POST("/products", handler.CreateProduct)
	r.GET("/products", handler.GetProducts)

	r.Run("localhost:8888")
}