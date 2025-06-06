package errors

import "github.com/go-kratos/kratos/v2/errors"

var (
	// 使用領域語言：即使錯誤集中定義，也可以使用領域語言命名錯誤。
	ErrDatabaseUnavailable = errors.ServiceUnavailable("DATABASE_UNAVAILABLE", "資料庫無法連線")
	ErrUserNotFound        = errors.NotFound("USER_NOT_FOUND", "使用者不存在")
	ErrInvalidParameter    = errors.BadRequest("INVALID_PARAM", "參數錯誤")
)
