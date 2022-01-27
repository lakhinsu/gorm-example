package routers

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/lakhinsu/gorm-example/controllers"
)

// Function to setup routers and router groups
func SetupRouters(app *gin.Engine) {
	v1 := app.Group("/v1")
	{
		v1.GET("/ping", controllers.Ping)
		v1.POST("/user", controllers.CreateUser)
		v1.GET("/user/:id", controllers.GetUser)
		v1.GET("/users", controllers.GetUsers)
		v1.PATCH("/user", controllers.UpdateUser)
		v1.DELETE("/user/:id", controllers.DeleteUser)
	}
	// Standalone route example
	// app.GET("/ping", controllers.Ping)
}
