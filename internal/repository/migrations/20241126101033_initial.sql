-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS account (
    id TEXT DEFAULT (hex(randomblob(8))) PRIMARY KEY NOT NULL,
    name TEXT NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS category (
    id TEXT DEFAULT (hex(randomblob(8))) PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    icon TEXT NOT NULL,
    color_hex TEXT NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
    account_id TEXT REFERENCES account(id) NOT NULL
) STRICT;

CREATE TABLE IF NOT EXISTS entry (
    id TEXT DEFAULT (hex(randomblob(8))) PRIMARY KEY NOT NULL,
    entry_type INTEGER NOT NULL, -- 0 - expense, 1 - income
    name TEXT NOT NULL,
    amount INTEGER NOT NULL, -- cents
    description TEXT,
    date TEXT DEFAULT CURRENT_DATE NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
    category_id TEXT REFERENCES category(id),
    account_id TEXT REFERENCES account(id) NOT NULL
) STRICT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE account;
DROP TABLE category;
DROP TABLE entry;
-- +goose StatementEnd
