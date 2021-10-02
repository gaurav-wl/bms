package server

import (
	"database/sql"
	"encoding/json"
	"github.com/gauravcoco/bms/crypto"
	"github.com/gauravcoco/bms/dbModels"
	"github.com/gauravcoco/bms/errors"
	"github.com/gauravcoco/bms/models"
	"github.com/gauravcoco/bms/utils"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

func (srv *Server) registerNewUser(resp http.ResponseWriter, req *http.Request) {
	var newUserRequest models.NewUserRequest

	err := json.NewDecoder(req.Body).Decode(&newUserRequest)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusBadRequest, "Error parsing request", "Error parsing new user request")
		return
	}

	newUser := &dbModels.User{
		Name:      newUserRequest.Name,
		Password:  crypto.HashAndSalt(newUserRequest.Password),
		Email:     newUserRequest.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	userID, err := srv.DBHelper.CreateUser(req.Context(), newUser)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error creating new user, Please try again", "Error parsing new user request")
		return
	}

	user, err := srv.DBHelper.GetUserByID(req.Context(), userID)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error creating new user, Please try again", "Error parsing new user request")
		return
	}

	token, err := getJWTToken(user.UUID, srv.Config.GetJWTKey())
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error getting user, Please try again", "Error getting user token")
		return
	}

	utils.EncodeJSONBody(resp, http.StatusCreated, map[string]interface{}{
		"token": token,
	})
}

func (srv *Server) loginUser(resp http.ResponseWriter, req *http.Request) {
	var loginReq models.LoginRequest

	err := json.NewDecoder(req.Body).Decode(&loginReq)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusBadRequest, "Error parsing request", "Error parsing login request")
		return
	}

	user := struct {
		UUID     string `json:"uuid"`
		Password string `json:"password"`
	}{}

	SQL := `SELECT 
				uuid,
				password
			FROM users 
			WHERE 
				email = $1 AND 
				archived_at IS NULL
				`

	err = srv.PSQL.DB().Get(&user, SQL, loginReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			bmsError.RespondClientErr(resp, req, err, http.StatusNotFound, "User Not exists")
			return
		}
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error getting user")
		return
	}

	if !crypto.ComparePasswords(user.Password, loginReq.Password) {
		bmsError.RespondClientErr(resp, req, err, http.StatusBadRequest, "Email or Password are incorrect")
		return
	}

	token, err := getJWTToken(user.UUID, srv.Config.GetJWTKey())
	if err != nil {
		bmsError.RespondGenericServerErr(resp, req, err, "Error generating jwt token")
		return
	}

	utils.EncodeJSONBody(resp, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}

func (srv *Server) getAllMovies(resp http.ResponseWriter, req *http.Request) {
	var movieSearchRequest models.MovieSearchRequest

	err := json.NewDecoder(req.Body).Decode(&movieSearchRequest)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusBadRequest, "Error parsing request", "Error parsing movies search request")
		return
	}

	movies, err := srv.DBHelper.GetAllMovies(req.Context(), movieSearchRequest)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Failed to get movies", "Error getting getAllMovies")
		return
	}

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"movies": srv.Converter.ToMovies(movies),
	})
}

func (srv *Server) getMovieDetails(resp http.ResponseWriter, req *http.Request) {
	movieID, err := strconv.Atoi(chi.URLParam(req, "movieID"))
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusBadRequest, "Error parsing movie id", "Error parsing movie id to int")
		return
	}

	movie, err := srv.DBHelper.GetMovieDetails(req.Context(), movieID)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error getting movie details", "Error getting movie details#" + strconv.Itoa(movieID))
		return
	}

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"movieDetails": srv.Converter.ToMovieDetails(*movie),
	})
}

func (srv *Server) getMovieShowDetails(resp http.ResponseWriter, req *http.Request) {
	movieID, err := strconv.Atoi(chi.URLParam(req, "movieID"))
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusBadRequest, "Error parsing movie id", "Error parsing movie id to int")
		return
	}

	showDetails, err := srv.DBHelper.GetMovieShowDetails(req.Context(), movieID)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error getting movie details", "Error getting movie details#" + strconv.Itoa(movieID))
		return
	}

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"movieShowDetails": srv.Converter.ToMovieShowDetails(*showDetails),
	})
}

func (srv *Server) getShowSeatsDetails(resp http.ResponseWriter, req *http.Request) {
	showID, err := strconv.Atoi(chi.URLParam(req, "showID"))
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusBadRequest, "Error parsing show id", "Error parsing show id to int")
		return
	}

	showSeatDetails, err := srv.DBHelper.GetShowSeats(req.Context(), showID)
	if err != nil {
		bmsError.RespondClientErr(resp, req, err, http.StatusInternalServerError, "Error getting movie details", "Error getting show seat details#" + strconv.Itoa(showID))
		return
	}

	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"movieShowSeatDetails": srv.Converter.ToSeatsDetails(showSeatDetails),
	})
}
