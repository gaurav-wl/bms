package middlewareProvider

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/bms/dbModels"
	bmsError "github.com/bms/errors"
	"github.com/bms/sql"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	authorization = "Authorization"
	bearerScheme  = "bearer"
	space         = " "
)

func corsOptions() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Access-Token", "importDate", "X-Client-Version", "Cache-Control", "Pragma"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

func authMiddleware(db *sqlx.DB, jwtKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.Split(r.Header.Get(authorization), space)
			if len(token) != 2 {
				bmsError.RespondClientErr(w, r, errors.New("token not Bearer"), http.StatusUnauthorized, "Invalid token")
				return
			}

			if strings.ToLower(token[0]) != bearerScheme {
				bmsError.RespondClientErr(w, r, errors.New("token not Bearer"), http.StatusUnauthorized, "Invalid token")
				return
			}

			claims := &dbModels.Claims{}

			tkn, err := jwt.ParseWithClaims(token[1], claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtKey), nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				logrus.Errorf("Error parsing jwt %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := getUserByUUID(db, claims.UUID)
			if err != nil {
				bmsError.RespondClientErr(w, r, err, http.StatusInternalServerError, "Error getting user context")
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), dbModels.UserContextKey, user))
			next.ServeHTTP(w, r)
		})
	}
}

func getUserByUUID(db *sqlx.DB, uuid string) (*dbModels.UserContext, error) {
	var user dbModels.UserContext
	err := db.Get(&user, sql.GetUserContextByUUIDSQL, uuid)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
