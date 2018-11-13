CREATE TABLE todos (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    is_completed BOOLEAN NOT NULL
)