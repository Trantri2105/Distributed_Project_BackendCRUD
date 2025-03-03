CREATE DATABASE users;

\c users

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name TEXT,
    email TEXT UNIQUE,
    password TEXT,
    phone TEXT,
    gender TEXT
);