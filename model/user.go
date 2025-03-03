package model

type User struct {
	Id       int    `db:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Phone    string `db:"phone" json:"phone"`
	Gender   string `db:"gender" json:"gender"`
}
