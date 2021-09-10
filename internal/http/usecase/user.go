package usecase

import (
	"net/http"
	"time"
	"twitter-clone/internal/http/gen"
	"twitter-clone/internal/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		db: db,
	}
}

type jwtCustomClaims struct {
	UID  string `json:"uid"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var signingKey = []byte("secret")

var Config = middleware.JWTConfig{
	Claims:     &jwtCustomClaims{},
	SigningKey: signingKey,
}

func (p *User) Signup(c echo.Context) error {
	// リクエストを取得
	user := new(gen.User)
	if err := c.Bind(user); err != nil {
		return sendError(c, http.StatusBadRequest, "Invalid format")
	}

	if user.Name == "" || user.Password == "" {
		return sendError(c, http.StatusBadRequest, "Invalid name or password")
	}

	u := new(repository.User)
	if p.db.Where("name = ?", user.Name).First(&u); u.ID != "" {
		return sendError(c, http.StatusConflict, "Name already exists")
	}

	p.db.Create(&repository.User{
		Name:     user.Name,
		Password: user.Password,
	})
	user.Password = ""

	return c.JSON(http.StatusOK, user)
}

func (p *User) Login(c echo.Context) error {
	// リクエストを取得
	user := new(gen.User)
	if err := c.Bind(user); err != nil {
		return sendError(c, http.StatusBadRequest, "Invalid format")
	}

	u := new(repository.User)
	p.db.Where("name = ?", user.Name).First(&u)
	if u.ID == "" || u.Password != user.Password {
		return sendError(c, http.StatusUnauthorized, "Invalid name or password")
	}

	claims := &jwtCustomClaims{
		u.ID,
		u.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(signingKey)
	if err != nil {
		return sendError(c, http.StatusBadRequest, "Sign in failed")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func userIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	uid := claims.UID
	return uid
}
