CREATE TABLE IF NOT EXISTS books (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    title varchar(255) NOT NULL,
    author varchar(255) NOT NULL,
    genre varchar(255) NOT NULL
);