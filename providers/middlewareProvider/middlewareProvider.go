package middlewareProvider

import (
	"github.com/gauravcoco/bms/providers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
)

type Middleware struct {
	DB     *sqlx.DB
	JWTKey string
}

func NewMiddleware(arDB *sqlx.DB, JWTKey string) providers.MiddlewareProvider {
	return Middleware{
		DB:     arDB,
		JWTKey: JWTKey,
	}
}

func (m Middleware) Default() chi.Middlewares {
	return chi.Chain(
		corsOptions().Handler,
		middleware.RequestID,
		middleware.RequestLogger(NewStructuredLogger()),
		middleware.Recoverer,
	)
}

func (m Middleware) AUTH() chi.Middlewares {
	return chi.Chain(
		authMiddleware(m.DB, m.JWTKey),
	)
}
