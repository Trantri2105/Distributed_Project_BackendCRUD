package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"reflect"
	"user_crud/model"
	"user_crud/repository"
	"user_crud/utils"
)

type UserService interface {
	RegisterUser(ctx context.Context, user model.User) error
	Login(ctx context.Context, email, password string) (string, error)
	UpdateUser(ctx context.Context, user model.User) error
	GetUserById(ctx context.Context, id int) (model.User, error)
}

type userService struct {
	userRepository repository.UserRepository
	jwtUtils       utils.JwtUtils
}

func (u *userService) RegisterUser(ctx context.Context, user model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	err = u.userRepository.InsertUser(ctx, user)
	return err
}

func (u *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("wrong password")
	}
	token, err := u.jwtUtils.CreateToken(user.Id)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userService) UpdateUser(ctx context.Context, user model.User) error {
	if !reflect.ValueOf(user.Password).IsZero() {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Teacher service, update teacher err :", err)
			return err
		}
		user.Password = string(hash)
	}
	err := u.userRepository.UpdateUser(ctx, user)
	return err
}

func (u *userService) GetUserById(ctx context.Context, id int) (model.User, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	return user, err

}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository, jwtUtils: utils.NewJwtUtils()}
}
