package utils

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"time"
)

type JwtUtils interface {
	CreateToken(userId int) (string, error)
	VerifyToken(tokenString string) (jwt.MapClaims, error)
}

const jwtTokenExpTime = 60 * time.Minute

type jwtUtils struct{}

func (*jwtUtils) CreateToken(userId int) (string, error) {
	expireTime := time.Now().Add(jwtTokenExpTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    expireTime,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		log.Println("Jwt service, create access token err :", err)
		return "", errors.New("internal server error")
	}
	return tokenString, nil
}

func (*jwtUtils) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("token is invalid")
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}

func NewJwtUtils() JwtUtils {
	return &jwtUtils{}
}
