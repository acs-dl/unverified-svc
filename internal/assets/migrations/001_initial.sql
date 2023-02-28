-- +migrate Up

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    module_id TEXT NOT NULL,
    username TEXT,
    phone TEXT,
    email TEXT,
    module TEXT NOT NULL,
    name TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (module_id, module)
);

CREATE INDEX IF NOT EXISTS users_username_idx ON users(username);
CREATE INDEX IF NOT EXISTS users_phone_idx ON users(phone);
CREATE INDEX IF NOT EXISTS users_email_idx ON users(email);
CREATE INDEX IF NOT EXISTS users_module_idx ON users(module);
CREATE INDEX IF NOT EXISTS users_moduleId_idx ON users(module_id);


-- +migrate Down

DROP TABLE IF EXISTS users;
DROP INDEX IF EXISTS users_username_idx;
DROP INDEX IF EXISTS users_phone_idx;
DROP INDEX IF EXISTS users_email_idx;
DROP INDEX IF EXISTS users_module_idx;
DROP INDEX IF EXISTS users_moduleId_idx;