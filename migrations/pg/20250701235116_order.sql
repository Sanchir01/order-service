-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS  orders(
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    track_number TEXT NOT NULL ,
    entry TEXT NOT NULL,
    locale TEXT NOT NULL,
    internal_signature TEXT NOT NULL ,
    customer_id TEXT NOT NULL,
    delivery_service TEXT NOT NULL,
    shardkey INT NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
