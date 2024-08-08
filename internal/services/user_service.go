package services

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id       int
	Username string
	Password string
}

type UserService struct {
	User      User
	UserStore *sql.DB
}

func NewUserService(u User, uStore *sql.DB) *UserService {
	return &UserService{User: u, UserStore: uStore}
}

func (us *UserService) CreateUser(user, pass string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(username, password) VALUES($1, $2)`

	_, err = us.UserStore.Exec(
		stmt,
		user,
		string(hashedPassword),
	)

	return err
}

func (us *UserService) CheckUser(username string) (User, error) {
	query := `SELECT id, username, password FROM users
		WHERE username = ?`

	stmt, err := us.UserStore.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.Username = username
	err = stmt.QueryRow(
		us.User.Username,
	).Scan(
		&us.User.id,
		&us.User.Username,
		&us.User.Password,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
}
