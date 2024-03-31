CREATE TABLE IF NOT EXISTS url(
                                  id SERIAL PRIMARY KEY,
                                  alias TEXT NOT NULL UNIQUE,
                                  url TEXT NOT NULL
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);

-- migrate -path migrations -database postgres://shortener:1234@localhost:5432/shortener?sslmode=disable up




