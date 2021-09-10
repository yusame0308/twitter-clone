package api

import (
	"twitter-clone/internal/http/gen"
	"twitter-clone/internal/http/usecase"
	"twitter-clone/internal/repository"

	om "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Run() {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// validator
	spec, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}
	e.Use(om.OapiRequestValidator(spec))

	// mysql connection
	dsn := "docker:docker@tcp(127.0.0.1:3306)/twitterCloneApi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	// Migrate the schema
	if err := db.AutoMigrate(&repository.User{}); err != nil {
		panic(err.Error())
	}
	api := e.Group("/api")
	api.Use(middleware.JWTWithConfig(usecase.Config))
	gen.RegisterHandlers(e, NewApi(db))
	e.Logger.Fatal(e.Start(":1232"))

	//e.POST("/signup", handler2.Signup)
	//e.POST("/login", handler.Login)
	//
	//api := e.Group("/api")
	//api.Use(middleware.JWTWithConfig(handler2.Config))
	//api.GET("/health", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "Get Health")
	//})
	//
	//return e
}
