package errors

import "github.com/go-kratos/kratos/v2/errors"

var (
	// 使用領域語言：即使錯誤集中定義，也可以使用領域語言命名錯誤。
	// Handlers
	ErrMissingRequiredField = errors.BadRequest("MISSING_REQUIRED_FIELD", "必填欄位不能為空")
	ErrInvalidParameter     = errors.BadRequest("INVALID_PARAM", "參數錯誤")
	ErrInvalidFieldFormat   = errors.BadRequest("INVALID_FIELD_FORMAT", "欄位格式無效")
	// Middleware
	ErrRequestIDGeneration = errors.InternalServer("REQUEST_ID_GENERATION_FAILED", "請求 ID 生成失敗")
	//
	ErrDatabaseUnavailable = errors.ServiceUnavailable("DATABASE_UNAVAILABLE", "資料庫無法連線")

	// Logger
	ErrLoggerNotFound    = errors.InternalServer("LOGGER_NOT_FOUND", "Logger not found in context")
	ErrInvalidLoggerType = errors.InternalServer("INVALID_LOGGER_TYPE", "Invalid logger type in context")

	// Repository
	// Wallet
	ErrWalletMnemonicRequired      = errors.InternalServer("WALLET_MNEMONIC_REQUIRED", "助記詞不能為空")
	ErrWalletInvalidDerivationPath = errors.InternalServer("WALLET_INVALID_PATH", "錢包衍生路徑格式錯誤")
	ErrWalletInvalidPathComponent  = errors.InternalServer("WALLET_INVALID_PATH_COMPONENT", "錢包衍生路徑分量無效")
	// Address 地址池
	ErrWalletInvalidAddressCount = errors.BadRequest("WALLET_INVALID_ADDRESS_COUNT", "地址數量必須為正數")
	ErrWalletUnsupportedChain    = errors.BadRequest("WALLET_UNSUPPORTED_CHAIN", "不支持的區塊鏈類型")
	ErrWalletAddressPoolUpdate   = errors.InternalServer("WALLET_ADDRESS_POOL_UPDATE", "地址池更新錯誤")
	// Logs 地址池更新紀錄
	ErrWalletAddressLogCreate = errors.InternalServer("WALLET_ADDRESS_LOG_CREATE", "地址池紀錄創建錯誤")
	// HDPath 相關錯誤
	ErrHDPathInvalidFormat    = errors.InternalServer("HDPATH_INVALID_FORMAT", "HD 路徑格式無效")
	ErrHDPathInvalidComponent = errors.InternalServer("HDPATH_INVALID_COMPONENT", "HD 路徑組件無效")
	ErrHDPathTooShort         = errors.InternalServer("HDPATH_TOO_SHORT", "HD 路徑太短，不符合 BIP44 標準")
	ErrHDPathInvalidPurpose   = errors.InternalServer("HDPATH_INVALID_PURPOSE", "HD 路徑的目的欄位無效")
	ErrHDPathMissingHardened  = errors.InternalServer("HDPATH_MISSING_HARDENED", "HD 路徑需要硬化派生")
	// Stream
	ErrStreamRedisCheckFailed      = errors.InternalServer("STREAM_REDIS_CHECK_FAILED", "檢查監聽地址時發生錯誤")
	ErrStreamAddressAlreadyWatched = errors.BadRequest("STREAM_ADDRESS_ALREADY_WATCHED", "該地址已經被監聽中")
	ErrStreamAddWatchFailed        = errors.InternalServer("STREAM_ADD_WATCH_FAILED", "添加監聽請求到流失敗")
	ErrStreamSetExpiryFailed       = errors.InternalServer("STREAM_SET_EXPIRY_FAILED", "設置監聽過期時間失敗")

	// Payment Order Validation
	ErrOrderNotUnique          = errors.InternalServer("ORDER_NOT_UNIQUE", "訂單筆數不為一筆")
	ErrOrderTokenNotFound      = errors.InternalServer("ORDER_TOKEN_NOT_FOUND", "找不到對應的代幣資訊")
	ErrOrderContractMismatch   = errors.InternalServer("ORDER_CONTRACT_MISMATCH", "合約地址不匹配")
	ErrOrderInsufficientAmount = errors.InternalServer("ORDER_INSUFFICIENT_AMOUNT", "付款金額不足")
)
