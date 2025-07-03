-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  payment(
                                        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                        order_uid UUID UNIQUE NOT NULL ,
                                        transaction UUID DEFAULT  uuid_generate_v4(),
                                        request_id TEXT NOT NULL DEFAULT '',
                                        currency TEXT NOT NULL,
                                        provider TEXT NOT NULL,
                                        amount INT NOT NULL,
                                        payment_dt INT NOT NULL,
                                        bank TEXT NOT NULL,
                                        delivery_cost INT NOT NULL ,
                                        goods_total INT NOT NULL ,
                                        custom_fee INT NOT NULL ,
                                        FOREIGN KEY (order_uid) REFERENCES orders(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
