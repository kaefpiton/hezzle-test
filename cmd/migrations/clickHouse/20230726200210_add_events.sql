-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS default.Events
(
    `Id`            UInt32,
    `CampaignId`    UInt64,
    `Name`          String,
    `Description`   String,
    `Priority`      UInt16,
    `Removed`       bool,
    `EventTime`     DateTime
)
engine = SummingMergeTree PARTITION BY toYYYYMM(EventTime)
ORDER BY (Id, CampaignId, Name)
SETTINGS index_granularity = 8192;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE  default.Events;
-- +goose StatementEnd
