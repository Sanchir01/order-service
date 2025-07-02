-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS  payment(
                                        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                        
                                        order_uid UUID UNIQUE NOT NULL ,
                                        currency TEXT NOT NULL,
                                        provider TEXT NOT NULL,
                                        amount INT NOT NULL,
                                        payment_dt TEXT NOT NULL,
                                        region TEXT NOT NULL,
                                        email TEXT NOT NULL,
                                        FOREIGN KEY (order_uid) REFERENCES orders(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
