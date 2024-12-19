-- +goose Up
-- +goose StatementBegin
INSERT INTO users VALUES (
  'e2a54cc6-814e-4fff-868b-8920c867e705',
  'admin@admin',
  'admin',
  NOW()::TIMESTAMP,
  NOW()::TIMESTAMP
);

-- +goose ENVSUB ON
insert into users_credentials values (
  'e2a54cc6-814e-4fff-868b-8920c867e705',
  '${ADMIN_PASSWORD}'
);
-- +goose ENVSUB OFF

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE uid = 'e2a54cc6-814e-4fff-868b-8920c867e705';
DELETE FROM users_credentials WHERE user_uid = 'e2a54cc6-814e-4fff-868b-8920c867e705';
-- +goose StatementEnd
