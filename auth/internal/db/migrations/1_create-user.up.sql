CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL -- bcrypt hash is 60 characters long
);

CREATE TABLE sessions
(
    token      VARCHAR(172) PRIMARY KEY, -- 128 byte token formatted in base64 is 172 characters long
    user_id    INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP,
    expires_at TIMESTAMP
);