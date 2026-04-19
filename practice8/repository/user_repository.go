package repository

// User is the domain model used across the repository layer.
type User struct {
	ID   int
	Name string
}

// UserRepository defines all data-access operations on User.
// Task 2: expanded with GetByEmail, UpdateUser, DeleteUser.
type UserRepository interface {
	GetUserByID(id int) (*User, error)
	CreateUser(user *User) error
	GetByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int) error
}
