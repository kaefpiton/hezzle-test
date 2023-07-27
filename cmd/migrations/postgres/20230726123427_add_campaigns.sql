-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS CAMPAIGNS (
    id serial NOT NULL,
    name VARCHAR(256) NOT NULL,

    PRIMARY KEY(id)
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE CAMPAIGNS;
-- +goose StatementEnd
