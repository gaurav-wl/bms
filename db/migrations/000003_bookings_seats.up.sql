ALTER TABLE bookings DROP COLUMN IF EXISTS seat_id;


CREATE TABLE if not exists bookings_seats
(
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES bookings(id),
    seat_id INTEGER NOT NULL REFERENCES movie_hall_seating(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS bookings_seats_unique ON bookings_seats(booking_id, seat_id);
