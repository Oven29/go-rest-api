-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_team (
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, team_id)
);

CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE IF NOT EXISTS pull_requests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author_id INTEGER REFERENCES users(id) ON DELETE RESTRICT,
    status pr_status NOT NULL DEFAULT 'OPEN'
);

CREATE TABLE IF NOT EXISTS pull_request_reviewer (
    pr_id INTEGER REFERENCES pull_requests(id) ON DELETE CASCADE,
    reviewer_id INTEGER REFERENCES users(id) ON DELETE RESTRICT,
    PRIMARY KEY (pr_id, reviewer_id)
);

CREATE INDEX idx_reviewer_id ON pull_request_reviewer (reviewer_id);


-- +goose Down
DROP TABLE IF EXISTS pull_request_reviewer;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS user_team;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS pr_status;
