-- +goose Up
-- +goose StatementBegin

insert into config(key, value) values
('mq_publish_messages_interval', '"15s"'),
('mq_publish_messages_batch_size', '10');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from config where key = 'mq_publish_messages_interval';
delete from config where key = 'mq_publish_messages_batch_size';
-- +goose StatementEnd
