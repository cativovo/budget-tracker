-- +goose Up
-- +goose StatementBegin
ALTER TABLE category
ADD account_id UUID NOT NULL REFERENCES account(id); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE category
DROP COLUMN account_id;
-- +goose StatementEnd
