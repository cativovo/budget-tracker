-- +goose Up
-- +goose StatementBegin

-- use the id, name, and email from the login provider
 CREATE TABLE user (
 	id TEXT NOT NULL PRIMARY KEY,
 	name TEXT NOT NULL,
 	email TEXT NOT NULL UNIQUE
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

CREATE INDEX idx_category_user_id ON category(user_id);

CREATE TABLE expense_group (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (hex(randomblob(8))),
	name TEXT NOT NULL,
	note TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE
);

CREATE INDEX idx_expense_group_user_id ON expense_group(user_id);

CREATE TABLE expense (
    id TEXT NOT NULL PRIMARY KEY DEFAULT (hex(randomblob(8))),
	name TEXT NOT NULL,
	amount INTEGER NOT NULL,
	note TEXT NOT NULL,
	date TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id TEXT NOT NULL REFERENCES user(id) ON DELETE CASCADE,
	category_id TEXT NOT NULL REFERENCES category(id) ON DELETE CASCADE,
	expense_group_id TEXT REFERENCES expense_group(id)
);

CREATE INDEX idx_expense_date ON expense(date);
CREATE INDEX idx_expense_user_id ON expense(user_id);
CREATE INDEX idx_expense_category_id ON expense(category_id);
CREATE INDEX idx_expense_expense_group_id ON expense(expense_group_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE user;

DROP TABLE category;
DROP INDEX idx_category_user_id;

DROP TABLE expense_group;
DROP INDEX idx_expense_group_user_id;

DROP TABLE expense;
DROP INDEX idx_expense_date;
DROP INDEX idx_expense_user_id;
DROP INDEX idx_expense_category_id;
DROP INDEX idx_expense_expense_group_id;

-- +goose StatementEnd
