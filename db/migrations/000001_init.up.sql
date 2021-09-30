CREATE TABLE if not exists users
(
    id SERIAL PRIMARY KEY,
    uuid uuid NOT NULL DEFAULT md5(random()::text || clock_timestamp()::text)::uuid,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    email_verified_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX uuid_unique_idx ON users(uuid);

CREATE TABLE if not exists users_last_location
(
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    last_waypoint point NOT NULL
);

CREATE TABLE if not exists state
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE if not exists city_states
(
    id SERIAL PRIMARY KEY,
    city_name TEXT NOT NULL,
    state_id integer NOT NULL REFERENCES state(id)
);

CREATE TYPE theater_status_enum AS ENUM (
  'active',
  'not_functional',
  'opening_soon'
);

CREATE TABLE if not exists theater
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    location point NOT NULL,
    city_id INTEGER NOT NULL REFERENCES city_states(id),
    address TEXT NOT NULL,
    status theater_status_enum NOT NULL DEFAULT 'active'
);

CREATE INDEX theater_status_idx ON theater(status);

CREATE TABLE if not exists movies
(
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    release_date date NOT NULL,
    duration_in_minutes INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE gender_enum AS ENUM (
  'male',
  'female',
  'other'
);


CREATE TYPE image_types AS ENUM (
  'movie_banner',
  'cast_image'
);

CREATE TABLE if not exists images
(
    id SERIAL PRIMARY KEY,
    bucket TEXT NOT NULL,
    path TEXT NOT NULL,
    type image_types NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    url TEXT NOT NULL
);

CREATE INDEX image_type_idx ON images(type);

CREATE TABLE if not exists movie_cast
(
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    cast_real_name TEXT NOT NULL,
    cast_reel_name TEXT NOT NULL,
    gender gender_enum NOT NULL,
    cast_image_url TEXT
);

CREATE TABLE if not exists movie_banners
(
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    banner_id INTEGER NOT NULL REFERENCES images(id),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TYPE movie_dimensions AS ENUM (
  '2D',
  '3D'
);


CREATE TABLE if not exists movie_dimension
(
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    dimension movie_dimensions NOT NULL
);

CREATE INDEX movie_dimension_idx ON movie_dimension(dimension);

CREATE TYPE movie_languages AS ENUM (
  'Hindi',
  'English',
  'Marathi',
  'Tamil',
  'Telugu'
);

CREATE TABLE if not exists movie_language
(
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    language movie_languages NOT NULL
);


CREATE TYPE movie_theater_status AS ENUM (
  'closed',
  'ongoing',
  'scheduled'
);

CREATE TABLE if not exists movie_theater_schedule
(
    id SERIAL PRIMARY KEY,
    theater_id INTEGER NOT NULL REFERENCES theater(id),
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    release_date TIMESTAMP WITH TIME ZONE NOT NULL,
    booking_start_time TIMESTAMP WITH TIME ZONE,
    status movie_theater_status NOT NULL DEFAULT 'scheduled',
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX movie_theater_schedule_unique_idx ON movie_theater_schedule(theater_id, movie_id) WHERE archived_at IS NULL;


CREATE TABLE if not exists movie_halls
(
    id SERIAL PRIMARY KEY,
    theater_id INTEGER NOT NULL REFERENCES theater(id),
    name TEXT NOT NULL,
    total_rows INTEGER NOT NULL,
    total_columns INTEGER NOT NULL,
    total_seats INTEGER NOT NULL
);

CREATE UNIQUE INDEX movie_halls_unique_idx ON movie_halls(theater_id, name);

CREATE TABLE if not exists show_timings
(
    id SERIAL PRIMARY KEY,
    theater_id INTEGER NOT NULL REFERENCES theater(id),
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    hall_id INTEGER NOT NULL REFERENCES movie_halls(id),
    show_start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    show_end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    dimension movie_dimensions NOT NULL,
    language movie_languages NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX show_timings_language_idx ON show_timings(language);
CREATE INDEX show_timings_dimension_idx ON show_timings(dimension);

CREATE TYPE hall_seats_status AS ENUM (
  'functional',
  'not_functional'
);

CREATE TABLE if not exists movie_hall_seating
(
    id SERIAL PRIMARY KEY,
    seat_code TEXT NOT NULL,
    hall_id INTEGER NOT NULL REFERENCES movie_halls(id),
    status hall_seats_status NOT NULL DEFAULT 'functional',
    row_number INTEGER NOT NULL,
    column_number INTEGER NOT NULL,
    is_recliner BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX movie_hall_seating_unique_idx ON movie_hall_seating(hall_id, seat_code);

CREATE TYPE booking_status AS ENUM (
  'confirmed',
  'cancelled'
);


CREATE TABLE if not exists bookings
(
    id SERIAL PRIMARY KEY,
    booking_pretty_id TEXT NOT NULL,
    user_id INTEGER NOT NULL REFERENCES users(id),
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    theater_id INTEGER NOT NULL REFERENCES theater(id),
    hall_id INTEGER NOT NULL REFERENCES movie_halls(id),
    seat_id INTEGER NOT NULL REFERENCES movie_hall_seating(id),
    status booking_status NOT NULL DEFAULT 'confirmed',
    show_id INTEGER NOT NULL REFERENCES show_timings(id),
    cancelled_reason TEXT,
    cancelled_at TIMESTAMP WITH TIME ZONE
);
