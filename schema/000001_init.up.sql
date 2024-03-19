CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    hash_password VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL,

    CONSTRAINT chk_role CHECK (role IN ('ordinary', 'admin'))
);

CREATE TABLE IF NOT EXISTS films
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(150),
    description VARCHAR(150),
    release_date DATE NOT NULL,
    rating SMALLINT NOT NULL,

    CONSTRAINT lim_title CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 150),
    CONSTRAINT lim_description CHECK (LENGTH(description) >= 0 AND LENGTH(description) <= 1000),
    CONSTRAINT lim_rating CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE IF NOT EXISTS actors
(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(150) NOT NULL,
    last_name VARCHAR(150) NOT NULL,
    surname VARCHAR(150),
    sex VARCHAR(255) NOT NULL,
    birthday DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS films_actors (
    film_id INT,
    actor_id INT,

    PRIMARY KEY (film_id, actor_id),
    FOREIGN KEY (film_id) REFERENCES films(id),
    FOREIGN KEY (actor_id) REFERENCES actors(id)
);