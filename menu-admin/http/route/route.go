package route

import (
	"menu_admin/http/handler"
	"menu_admin/http/middleware"

	"github.com/labstack/echo/v4"
)

func Product(e *echo.Echo) {
	g := e.Group("api/v1/product")
	g.POST("/", middleware.AuthorizeIsLogin(handler.CreateProduct))
	g.PUT("/:id", middleware.AuthorizeIsLogin(handler.UpdateProduct))
	g.DELETE("/:id", middleware.AuthorizeIsLogin(handler.DeleteProduct))
}

func User(e *echo.Echo) {
	g := e.Group("api/v1/user")
	g.GET("/", middleware.AuthorizeIsLogin(handler.GetAllUser))
	g.PATCH("/password/", middleware.AuthorizeIsLogin(handler.UpdateUserPassword))
	g.POST("/", handler.CreateUser)
	g.POST("/login/", handler.LoginUser)
	g.GET("/login/", middleware.AuthorizeIsLogin(handler.GetMyUser))
	g.GET("/logout/", middleware.DeleteSession)
}

func All(e *echo.Echo) {
	Product(e)
	User(e)
}
