-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    module_id BIGINT NOT NULL,
    username TEXT UNIQUE,
    phone TEXT UNIQUE,
    email TEXT UNIQUE,
    module TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (module_id, module)
);

CREATE INDEX IF NOT EXISTS users_idx ON users(username, phone, email, module, module_id);


-- +migrate Down

DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS users_idx;