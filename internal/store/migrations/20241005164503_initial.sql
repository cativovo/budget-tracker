-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS category (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	icon VARCHAR(20) NOT NULL,
	color_hex VARCHAR(6) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS account (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS expense (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	amount DECIMAL NOT NULL,
	description TEXT,
	date DATE DEFAULT CURRENT_DATE,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	category_id UUID REFERENCES category(id),
	account_id UUID NOT NULL REFERENCES account(id)
);

CREATE TABLE IF NOT EXISTS income (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	amount DECIMAL NOT NULL,
	description TEXT,
	date DATE DEFAULT CURRENT_DATE,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	category_id UUID REFERENCES category(id),
	account_id UUID NOT NULL REFERENCES account(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE category;
DROP TABLE expense;
DROP TABLE income;
DROP TABLE account;
-- +goose StatementEnd
