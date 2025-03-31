use card;

CREATE TABLE cards (
    id SERIAL PRIMARY KEY,
    word TEXT NOT NULL,
    meaning TEXT NOT NULL
);