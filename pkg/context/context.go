package context

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// MyContext
// echo.Contextをラップする構造体
type MyContext struct {
	echo.Context
}

// AuthBind
// Auth認証とBindを合わせたメソッド
func (c *MyContext) AuthBind(i interface{}) error {
	fmt.Println("Auth Bind")
	if err := c.Bind(i); err != nil {
		return err
	}
	return nil
}
