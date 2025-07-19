package types

type User struct {
	ID       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Phone    string `db:"phone"`
	Password string `db:"password"`
}
