package api

import (
	"context"
	"fmt"
	"twitter-clone/internal/http/gen"
	"twitter-clone/internal/http/usecase"
	"twitter-clone/internal/repository"

	"github.com/getkin/kin-openapi/openapi3filter"

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
	validatorOptions := &om.Options{}
	validatorOptions.Options.AuthenticationFunc = func(c context.Context, input *openapi3filter.AuthenticationInput) error {
		fmt.Println(">>>> INSIDE AuthenticationFunc")
		e.Use(middleware.JWTWithConfig(usecase.Config))
		return nil
	}
	e.Use(om.OapiRequestValidatorWithOptions(spec, validatorOptions))

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
	if err := db.AutoMigrate(&repository.Tweet{}); err != nil {
		panic(err.Error())
	}
	gen.RegisterHandlers(e, NewApi(db))
	e.Logger.Fatal(e.Start(":1232"))
}
