package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"reflect"
	"strings"
	"user_crud/model"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user model.User) error
	UpdateUser(ctx context.Context, user model.User) error
	GetUserById(ctx context.Context, id int) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	row := u.db.QueryRowxContext(ctx, query, email)
	var user model.User
	err := row.StructScan(&user)
	return user, err
}

func (u *userRepository) InsertUser(ctx context.Context, user model.User) error {
	query := `INSERT INTO users (name, email, password, phone, gender) VALUES 
            (:name, :email, :password, :phone, :gender)`

	_, err := u.db.NamedExecContext(ctx, query, user)
	return err
}

func (u *userRepository) UpdateUser(ctx context.Context, user model.User) error {
	var updateFields []string
	t := reflect.TypeOf(user)
	v := reflect.ValueOf(user)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		if !value.IsZero() {
			updateFields = append(updateFields, field.Tag.Get("db")+" = :"+field.Tag.Get("db"))
		}
	}
	if len(updateFields) == 0 {
		return nil
	}
	query := `UPDATE users SET ` + strings.Join(updateFields, ",") + ` WHERE id = :id`

	_, err := u.db.NamedExecContext(ctx, query, user)
	return err
}
func (u *userRepository) GetUserById(ctx context.Context, id int) (model.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	row := u.db.QueryRowxContext(ctx, query, id)
	var us model.User
	err := row.StructScan(&us)
	return us, err
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}
