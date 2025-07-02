-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  items(
                                       id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                       track_number TEXT NOT NULL,
                                       price INT NOT NULL,
                                       name TEXT NOT NULL ,
                                       sale TEXT NOT NULL,
                                       size INT NOT NULL,
                                       total_price INT NOT NULL,
                                       nm_id INT NOT NULL,
                                       brand TEXT NOT NULL,
                                       status int NOT NULL
);
CREATE TABLE order_items (
                             id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                             order_uid UUID UNIQUE NOT NULL ,
                             item_id UUID NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             FOREIGN KEY (order_uid) REFERENCES orders(id) ON DELETE CASCADE,
                             FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE RESTRICT

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
