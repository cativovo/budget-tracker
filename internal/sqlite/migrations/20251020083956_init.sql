-- +goose Up
-- +goose StatementBegin

-- use the id, name, and email from the login provider
 CREATE TABLE user (
 	id TEXT NOT NULL PRIMARY KEY,
 	name TEXT NOT NULL,
 	email TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
 );

CREATE TABLE category (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (hex(randomblob(8))),
	name TEXT NOT NULL,
	color TEXT NOT NULL,
	icon TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE
);

CREATE TABLE expense_group (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (hex(randomblob(8))),
	name TEXT NOT NULL,
	note TEXT NOT NULL,
	date DATE NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
	category_id TEXT NOT NULL REFERENCES category(id) ON DELETE CASCADE
);

CREATE TABLE expense (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (hex(randomblob(8))),
	name TEXT NOT NULL,
	amount INTEGER NOT NULL,
	note TEXT NOT NULL,
	date DATE NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
	category_id TEXT NOT NULL REFERENCES category(id) ON DELETE CASCADE,
	expense_group_id TEXT REFERENCES expense_group(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE user;
DROP TABLE category;
DROP TABLE expense_group;
DROP TABLE expense;

-- +goose StatementEnd
