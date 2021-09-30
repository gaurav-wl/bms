package dbModels

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	UserContextKey = "user_context"
)

type UserContext struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
}

type User struct {
	ID        int       `db:"id"`
	UUID      string    `db:"uuid"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}

type Claims struct {
	UUID string `json:"key"`
	jwt.StandardClaims
}
