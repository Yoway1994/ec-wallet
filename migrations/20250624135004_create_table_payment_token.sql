-- +goose Up
-- +goose StatementBegin
CREATE TABLE payment_tokens (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(64) NOT NULL,
    chain VARCHAR(64) NOT NULL,
    contract_address VARCHAR(1024),
    decimals INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT uk_token_identity CHECK (
        (contract_address IS NOT NULL) OR
        (contract_address IS NULL AND symbol IS NOT NULL AND chain IS NOT NULL)
    )
);

CREATE UNIQUE INDEX uq_payment_tokens_token ON payment_tokens(chain, contract_address)
WHERE contract_address IS NOT NULL;

CREATE UNIQUE INDEX uq_payment_tokens_native ON payment_tokens(chain, symbol)
WHERE contract_address IS NULL;


CREATE INDEX idx_payment_tokens_symbol ON payment_tokens(symbol);
CREATE INDEX idx_payment_tokens_chain ON payment_tokens(chain);
CREATE INDEX idx_payment_tokens_contract_address ON payment_tokens(contract_address);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment_tokens;
-- +goose StatementEnd
