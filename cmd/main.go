package main

import (
	"context"
	"fmt"
	"twitter-clone/infra/mysql/repository"
	api "twitter-clone/internal/http"
	"twitter-clone/internal/http/gen"
	mc "twitter-clone/pkg/context"
	mv "twitter-clone/pkg/validator"

	"github.com/getkin/kin-openapi/openapi3filter"
	"gopkg.in/go-playground/validator.v9"

	om "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	// echo validator
	e.Validator = &mv.Validator{Validator: validator.New()}

	// openAPI validator
	spec, err := gen.GetSwagger()
	if err != nil {
		panic(err)
	}
	validatorOptions := &om.Options{}
	validatorOptions.Options.AuthenticationFunc = func(c context.Context, input *openapi3filter.AuthenticationInput) error {
		fmt.Println(">>>> INSIDE AuthenticationFunc")
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

	// echo.Contextをラップするためにmiddlewareとして登録
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&mc.MyContext{
				Context: c,
				DB:      db,
				Auth: &mc.Auth{
					Name:     "name",
					Password: "pass",
				},
			})
		}
	})

	gen.RegisterHandlers(e, api.NewApi())
	e.Logger.Fatal(e.Start(":1232"))
}
