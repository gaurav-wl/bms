package server

import (
	"github.com/go-chi/chi"
)

func (srv *Server) InjectRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(srv.Middlewares.Default()...)
	r.Get(`/health`, srv.healthCheck)
	r.Route("/api", func(api chi.Router) {
		api.Post("/movies", srv.getAllMovies)
		api.Get("/movie/{movieID}", srv.getMovieDetails)
		api.Get("/movie_shows/{movieID}", srv.getMovieShowDetails)

		api.Route("/user", func(public chi.Router) {
			public.Post("/register", srv.registerNewUser)
			public.Post("/login", srv.loginUser)

			public.Route("/", func(user chi.Router) {
				user.Use(srv.Middlewares.AUTH()...)
				user.Get("/info", srv.userInfo)
				user.Post("/book", srv.book)
				public.Get("/seats/{showID}", srv.getShowSeatsDetails)
				user.Get("/bookings", srv.getAllBookings)
			})
		})
	})
	return r
}
