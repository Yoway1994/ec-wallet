-- +goose Up
-- +goose StatementBegin
-- 地址池表
CREATE TABLE wallet_address_pools (
    id SERIAL PRIMARY KEY,
    address VARCHAR(128) NOT NULL,
    chain VARCHAR(20) NOT NULL,
    
    -- 地址資訊
    path VARCHAR(100) NOT NULL,
    index INT NOT NULL,
    
    -- 狀態管理
    current_status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
    
    -- 基本時間戳
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- 唯一約束
    UNIQUE(address, chain)
);

-- 為表和欄位添加註釋
COMMENT ON TABLE wallet_address_pools IS '錢包地址資源池表';
COMMENT ON COLUMN wallet_address_pools.chain IS '鏈標識符 (ETH, BTC, SOL 等)';
COMMENT ON COLUMN wallet_address_pools.path IS '衍生路徑';
COMMENT ON COLUMN wallet_address_pools.index IS '衍生路徑中的索引值';
COMMENT ON COLUMN wallet_address_pools.current_status IS '當前狀態: AVAILABLE, RESERVED, BLACKLISTED';

-- 地址使用日誌表
CREATE TABLE wallet_address_logs (
    id SERIAL PRIMARY KEY,
    address_id INT NOT NULL REFERENCES wallet_address_pools(id),
    
    -- 操作和狀態
    operation VARCHAR(50) NOT NULL,
    status_after VARCHAR(20) NOT NULL,
    status_before VARCHAR(20),
    
    -- 時間資訊
    operation_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_until TIMESTAMP,
    
    -- 業務關聯資訊
    order_id VARCHAR(100),
    user_id VARCHAR(100)
);

-- 為日誌表和欄位添加註釋
COMMENT ON TABLE wallet_address_logs IS '錢包地址使用和狀態變更日誌表';
COMMENT ON COLUMN wallet_address_logs.operation IS '操作類型: ACQUIRE, RELEASE, BLACKLIST 等';
COMMENT ON COLUMN wallet_address_logs.status_after IS '操作後狀態: AVAILABLE, RESERVED, BLACKLISTED';
COMMENT ON COLUMN wallet_address_logs.status_before IS '操作前狀態';
COMMENT ON COLUMN wallet_address_logs.operation_at IS '操作執行時間';
COMMENT ON COLUMN wallet_address_logs.valid_until IS '如果是預約，有效期至';
COMMENT ON COLUMN wallet_address_logs.order_id IS '訂單ID';
COMMENT ON COLUMN wallet_address_logs.user_id IS '用戶ID';

-- 建立索引
CREATE INDEX idx_wallet_pools_status_chain ON wallet_address_pools(current_status, chain);
CREATE INDEX idx_wallet_logs_address_id ON wallet_address_logs(address_id);
CREATE INDEX idx_wallet_logs_operation_at ON wallet_address_logs(operation_at);
CREATE INDEX idx_wallet_logs_order_id ON wallet_address_logs(order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- 刪除索引
DROP INDEX IF EXISTS idx_wallet_logs_order_id;
DROP INDEX IF EXISTS idx_wallet_logs_operation_at;
DROP INDEX IF EXISTS idx_wallet_logs_address_id;
DROP INDEX IF EXISTS idx_wallet_pools_status_chain;

-- 刪除表格
DROP TABLE IF EXISTS wallet_address_logs;
DROP TABLE IF EXISTS wallet_address_pools;
-- +goose StatementEnd
