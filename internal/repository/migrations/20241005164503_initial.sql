-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS account (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS category (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	icon VARCHAR(20) NOT NULL,
	color_hex VARCHAR(6) NOT NULL,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	account_id UUID NOT NULL REFERENCES account(id)
);

CREATE TYPE transaction_type AS ENUM ('expense', 'income');

CREATE TABLE IF NOT EXISTS transaction (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	transaction_type transaction_type NOT NULL,
	name VARCHAR(255) NOT NULL,
	amount BIGINT NOT NULL, -- cents
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
DROP TABLE account;
DROP TABLE transaction;
-- +goose StatementEnd
