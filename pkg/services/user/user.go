package services

import "database/sql"

type User struct {
	ID    string
	Name  string
	Email string
	Role  string
	private
}

type private struct {
	Password string
	Cart     interface{}
}

type UserRepo interface {
	GetUserByID(id string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id string) error
}

func (u *UserService) GetUserByID(id string) (*User, error) {
	user := &User{}
	row := u.DB.QueryRow("SELECT * FROM users WHERE ID = $1", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) CreateUser(user *User) error {
	_, err := u.DB.Exec("INSERT INTO users (name, email, password, cart, role) VALUES ($1, $2, $3, $4, $5)", user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateUser(user *User) error {
	_, err := u.DB.Exec("UPDATE users SET name = $2, email = $3, role = $4, password = $5 WHERE id = $1", user.ID, user.Name, user.Email, user.Role, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) DeleteUser(id string) error {
	_, err := u.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}
