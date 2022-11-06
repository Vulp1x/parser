-- +goose Up
-- +goose StatementBegin
ALTER TABLE datasets
    ALTER COLUMN phone_code SET DATA TYPE integer;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin


-- +goose StatementEnd
