package server

import (
	"github.com/go-chi/chi"
)

func (srv *Server) InjectRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(srv.Middlewares.Default()...)
	r.Get(`/health`, srv.healthCheck)
	r.Route("/api", func(api chi.Router) {
		api.Route("/user", func(public chi.Router) {
			public.Post("/register", srv.registerNewUser)
			public.Post("/login", srv.loginUser)

			public.Post("/movies", srv.getAllMovies)
			public.Get("/movie/{id}", srv.getMovieDetails)
			public.Get("/movie_shows/{id}", srv.getMovieShowDetails)

			public.Route("/", func(user chi.Router) {
				user.Use(srv.Middlewares.AUTH()...)
				user.Get("/info", srv.userInfo)
				user.Post("/book", srv.userInfo)
				user.Get("/bookings", srv.userInfo)
			})
		})
	})
	return r
}
