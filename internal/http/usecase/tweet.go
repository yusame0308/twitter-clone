package usecase

import (
	"net/http"
	"strings"
	"twitter-clone/pkg/context"

	"gorm.io/gorm"
)

type Tweet struct {
	db *gorm.DB
}

func NewTweet(db *gorm.DB) *Tweet {
	return &Tweet{
		db: db,
	}
}

func (p *Tweet) CreateTweet(c *context.MyContext) error {
	//// リクエストを取得
	//tweet := new(gen.Tweet)
	//if err := c.Bind(tweet); err != nil {
	//	return sendError(c, http.StatusBadRequest, "Invalid format")
	//}
	//userID := userIDFromToken(c)
	//
	//tx := p.db.Create(&repository.Tweet{
	//	Text:      tweet.Text,
	//	UserID:    userID,
	//})
	//if tx.Error != nil {
	//	return sendError(c, http.StatusBadRequest, tx.Error.Error())
	//}

	//TODO 最初にlogin/signupするのは通るけど、tweetをした後だと通らなくなるから修正

	auth := c.Request().Header.Get("Authorization")
	if !strings.HasPrefix(strings.ToLower(auth), "bearer") {
		return c.String(http.StatusBadRequest, "hasPrefix")
	}

	values := strings.Split(auth, " ")
	if len(values) < 2 {
		return c.String(http.StatusBadRequest, "split")
	}

	token := values[1]

	return c.String(http.StatusCreated, token)
}
