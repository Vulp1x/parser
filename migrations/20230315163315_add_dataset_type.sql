-- +goose Up
-- +goose StatementBegin

alter table full_targets
    drop column bio_links,
    add column bio_links text not null default '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE dataset_type;
ALTER TABLE datasets
    DROP COLUMN type;
-- +goose StatementEnd
