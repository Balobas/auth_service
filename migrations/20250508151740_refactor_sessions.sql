-- +goose Up
-- +goose StatementBegin

alter table sessions add column tokens_issued_at integer;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table sessions drop column tokens_issued_at;
-- +goose StatementEnd
