ALTER TABLE movie_cast DROP COLUMN IF EXISTS cast_image_url;

ALTER TABLE movie_cast ADD COLUMN IF NOT EXISTS image_id INTEGER REFERENCES images(id);
