-- Active: 1715707344079@@127.0.0.1@5432@crypto-up

CREATE TABLE users (
    id           SERIAL PRIMARY KEY,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    email        VARCHAR(255) UNIQUE NOT NULL,
    phone        VARCHAR(20) UNIQUE,
    login        VARCHAR(255),
    display_name VARCHAR(255),
    image_url    VARCHAR(255),
    password     VARCHAR(255) NOT NULL,
    token        VARCHAR(255)
);
