package api

import (
	"twitter-clone/internal/http/gen"
	"twitter-clone/internal/http/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Api struct {
	user  *usecase.User
	tweet *usecase.Tweet
}

func NewApi(db *gorm.DB) *Api {
	return &Api{user: usecase.NewUser(db), tweet: usecase.NewTweet(db)}
}

var _ gen.ServerInterface = (*Api)(nil)

func (p Api) Login(ctx echo.Context) error {
	return p.user.Login(ctx)
}

func (p Api) Signup(ctx echo.Context) error {
	return p.user.Signup(ctx)
}

func (p Api) CreateTweet(ctx echo.Context) error {
	return p.tweet.CreateTweet(ctx)
}
