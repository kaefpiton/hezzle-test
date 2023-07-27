-- +goose Up
-- +goose StatementBegin
INSERT INTO CAMPAIGNS (name) VALUES('Первая запись');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE CAMPAIGNS;
-- +goose StatementEnd
