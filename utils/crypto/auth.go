package crypto

import (
	"serverless-todo-golang/utils/logger"
	"serverless-todo-golang/utils/middleware"
	"time"
)

type UserTokenData struct {
	ID        string
	CreatedAt time.Time
}

func (u *UserTokenData) TimeStamp() {
	u.CreatedAt = time.Now()
}

func GenerateAuthToken(tokenData UserTokenData) (string, error) {
	tokenData.TimeStamp()
	token, err := middleware.GenerateToken(&tokenData)
	if err != nil {
		logger.GetLog().Error("ERROR : ", err.Error())
		return "", err
	}
	return token, nil
}
