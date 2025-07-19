package types

type UserRole string

const (
	CustomerRole UserRole = "user"
	AdminRole    UserRole = "admin"
)

type User struct {
	ID       string   `db:"id"`
	Name     string   `db:"name"`
	Email    string   `db:"email"`
	Phone    string   `db:"phone"`
	Role     UserRole `db:"role"`
	Password string   `db:"password"`
}
