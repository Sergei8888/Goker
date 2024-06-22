CREATE TABLE users {
    id SERIAL PRIMARY KEY INT,
}

CREATE TABLE sessions {
    id SERIAL PRIMARY KEY INT,
    FOREIGN KEY user_id INT REFERENCES users(id),
    token VARCHAR(172), -- 128 byte token formatted in base64 is 172 characters long
    created_at TIMESTAMP,
    expires_at TIMESTAMP,
}