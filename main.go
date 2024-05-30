package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	database "github.com/mayank12gt/ProgressTracker/db"
	"github.com/mayank12gt/ProgressTracker/middleware"
	"github.com/mayank12gt/ProgressTracker/routes"
)

func main() {
	database.ConnectDB()
	e := echo.New()
	setupRoutes(e)
	e.Logger.Fatal(e.Start(":3000"))

}

func setupRoutes(e *echo.Echo) {

	//e.Use(middleware.Logger)
	e.GET("/user/profile", routes.Profile, middleware.VerifyJWT)
	e.GET("/", sayHello)
	e.POST("/user/signUp", routes.SignUp)
	e.POST("/user/signIn", routes.SignIn)
	e.GET("/user/signOut", routes.SignOut)

	list := e.Group("/list")
	list.Use(middleware.VerifyJWT)
	list.GET("/all", routes.GetAllLists)
	list.POST("/create", routes.CreateList)
	list.GET("/:id", routes.GetList)
	list.POST("/:id/update", routes.UpdateList)
	list.DELETE("/:id/delete", routes.DeleteList)

	item := list.Group("/:id/item")
	item.GET("/all", routes.GetAllItems)
	item.POST("/create", routes.CreateItem)
	item.DELETE("/:itemId/delete", routes.DeleteItem)
	item.GET("/:itemId", routes.GetItem)
	item.POST("/:itemId/update", routes.UpdateItem)
	item.GET("/:itemId/toggleComplete", routes.ToggleItemCompleted)

}

func sayHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello")
}
