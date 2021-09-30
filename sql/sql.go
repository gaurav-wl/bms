package sql

//language=SQL
const (

	CreateNewUserSQL = `INSERT INTO users
						(name, email, password)
						 VALUES ($1, $2, $3)
						 RETURNING id`

	GetUserByIdSQL = `SELECT
						id,  
						name, 
						email, 
						password
					  FROM
						 users 
					  WHERE
						  id = $1 AND 
						  archived_at IS NULL`

	GetUserInfoByIdSQL = `SELECT
							id,  
							name, 
							email 
						  FROM
							 users 
						  WHERE
							  id = $1 AND 
							  archived_at IS NULL`

	GetUserContextByUUIDSQL = `  SELECT
								    id,  
									name, 
									email 
								  FROM
									 users 
								  WHERE
									  uuid = $1 AND 
									  archived_at IS NULL`


	MovieShowDetailsSQL = `SELECT 
								DISTINCT movies.id AS movie_id,
								show_timings.id AS show_id,
								show_timings.show_start_time,
								show_timings.show_end_time,
								theater.id AS theater_id,
								theater.name AS theater_name,
								show_timings.dimension, 
								show_timings.language
							FROM movies
							JOIN movie_theater_schedule mts ON mts.movie_id = movies.id
							JOIN theater ON theater.id = mts.theater_id
							JOIN show_timings ON show_timings.movie_id = movies.id AND show_timings.theater_id = theater.id
							WHERE
								movies.id = $1 AND
								mts.booking_start_time > now() AND 
								mts.status = 'ongoing' AND 
								mts.archived_at IS NULL AND 
								show_timings.show_start_time > now() + '15 minutes'::INTERVAL AND 
								show_timings.archived_at IS NULL AND 
								show_timings.show_end_time::DATE < now()::DATE + '10 DAYS'::INTERVAL				
`

	MovieDetailsSQL  = `SELECT 
							movies.id,
							movies.name,
							movies.release_date,
							movies.duration_in_minutes,
							movie_dimensions.dimensions,
							movie_languages.languages
						FROM movies
						JOIN LATERAL (
								SELECT ARRAY_AGG(md.dimension) AS dimensions
								FROM movie_dimension md 
								WHERE md.movie_id = movies.id
						) movie_dimensions ON TRUE
						JOIN LATERAL (
								SELECT ARRAY_AGG(ml.language) AS languages
								FROM movie_languages ml 
								WHERE ml.movie_id = movies.id
						) movie_languages ON TRUE
						WHERE 
							movies.id = $1 AND 
							archived_at IS NULL
							`


)
