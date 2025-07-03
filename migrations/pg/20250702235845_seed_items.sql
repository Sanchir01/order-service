-- +goose Up
-- +goose StatementBegin
INSERT INTO items (track_number, price, name, sale, size, total_price, nm_id, brand, status) VALUES
                                                                                                 ('TRACK001', 100, 'Item One', '10%', 1, 90, 101, 'BrandA', 1),
                                                                                                 ('TRACK002', 200, 'Item Two', '15%', 2, 170, 102, 'BrandB', 1),
                                                                                                 ('TRACK003', 150, 'Item Three', '5%', 3, 142, 103, 'BrandC', 1),
                                                                                                 ('TRACK004', 300, 'Item Four', '20%', 4, 240, 104, 'BrandD', 1),
                                                                                                 ('TRACK005', 250, 'Item Five', '0%', 5, 250, 105, 'BrandE', 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
