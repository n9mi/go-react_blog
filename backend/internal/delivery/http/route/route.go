package route

import (
	"backend/internal/delivery/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var parentRoute = "/api"

type RouteConfig struct {
	App            *echo.Echo
	PostController *http.PostController
	UserController *http.UserController
	AuthMiddleware echo.MiddlewareFunc
}

func (r *RouteConfig) Setup() {
	r.SetupCommon()
	r.SetupGuestRoute()
	r.SetupAuthRoute()
	r.SetupUserRoute()
	r.SetupAdminRoute()
}

func (r *RouteConfig) SetupCommon() {
	r.App.Use(middleware.Logger())
	r.App.Use(middleware.Recover())
	r.App.Use(middleware.RemoveTrailingSlash())
	r.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
}

func (r *RouteConfig) SetupGuestRoute() {
	routeGroup := "/posts"

	g := r.App.Group(parentRoute + routeGroup)
	g.GET("", r.PostController.GetAll)
	g.GET("/:id", r.PostController.GetByID)
}

func (r *RouteConfig) SetupAuthRoute() {
	routeGroup := "/auth"

	g := r.App.Group(parentRoute + routeGroup)

	g.POST("/register", r.UserController.Register)
	g.POST("/login", r.UserController.Login)
	g.POST("/logout", r.UserController.Logout, r.AuthMiddleware)
	g.POST("/refresh", r.UserController.Refresh)
}

func (r *RouteConfig) SetupUserRoute() {
	routeGroup := "/user"

	g := r.App.Group(parentRoute + routeGroup)

	g.GET("/current", r.UserController.Current, r.AuthMiddleware)
}

func (r *RouteConfig) SetupAdminRoute() {
	routeGroup := "/admin"

	g := r.App.Group(parentRoute + routeGroup)
	g.Use(r.AuthMiddleware)

	g.POST("/posts", r.PostController.Create)
	g.PUT("/posts/:id", r.PostController.Update)
	g.DELETE("/posts/:id", r.PostController.Delete)
}
