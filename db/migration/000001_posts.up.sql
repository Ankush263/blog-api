CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR,
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- migrate create -ext sql -dir db/migration -seq posts -format
-- migrate -path ./db/migration -database "postgres://postgres:postgres@localhost/blogapi?sslmode=disable" up 1