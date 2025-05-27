CREATE TABLE IF NOT EXISTS movies (
    movie_id      SERIAL PRIMARY KEY,
    title         VARCHAR(255) NOT NULL,
    year          INTEGER,
    plot          TEXT,
    genre         VARCHAR(255),
    imdbid        VARCHAR(20) UNIQUE,
    actors        TEXT
);