-- +goose Up
-- +goose StatementBegin
CREATE TABLE payment_orders (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(64) NOT NULL,
    address VARCHAR(128) NOT NULL,
    chain VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    amount_usd DECIMAL(38, 18) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    tx_hash VARCHAR(128),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    expire_time TIMESTAMP WITH TIME ZONE NOT NULL,
    paid_at TIMESTAMP WITH TIME ZONE,
    
    UNIQUE(order_id),
    UNIQUE(address, chain)
);

CREATE INDEX idx_payment_orders_order_id ON payment_orders(order_id);
CREATE INDEX idx_payment_orders_address_chain ON payment_orders(address, chain);
CREATE INDEX idx_payment_orders_status ON payment_orders(status);
CREATE INDEX idx_payment_orders_created_at ON payment_orders(created_at);
CREATE INDEX idx_payment_orders_expire_time ON payment_orders(expire_time);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_orders;
-- +goose StatementEnd
