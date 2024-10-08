-- +goose Up
-- +goose StatementBegin
DROP TABLE expense;
DROP TABLE income;

CREATE TABLE IF NOT EXISTS transaction (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	transaction_type SMALLINT NOT NULL, -- 0 - expense, 1 - income; enum is hassle to setup
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
DROP TABLE expense;
DROP TABLE income;
-- +goose StatementEnd
