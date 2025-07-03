-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS delivery (
                                        id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
                                        order_uid UUID UNIQUE NOT NULL,
                                        name TEXT NOT NULL,
                                        phone TEXT NOT NULL,
                                        zip INT NOT NULL ,
                                        city TEXT NOT NULL,
                                        address TEXT NOT NULL,
                                        region TEXT NOT NULL,
                                        email TEXT NOT NULL,
                                        FOREIGN KEY (order_uid) REFERENCES orders(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
