-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ITEMS (
    id serial NOT NULL,
    campaign_id int NOT NULL,
    name VARCHAR(256) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    priority INT NOT NULL DEFAULT 1,
    removed bool NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY(id, campaign_id)
    );

CREATE INDEX ON ITEMS(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ITEMS;
-- +goose StatementEnd
