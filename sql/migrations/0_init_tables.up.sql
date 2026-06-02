
CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    username VARCHAR(255),
    subscribed BOOLEAN,
    language_code CHAR(2),
    date_add INTEGER
);

CREATE TABLE IF NOT EXISTS games (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    link VARCHAR(255),
    image_link VARCHAR(255),
    date_add INTEGER
);

