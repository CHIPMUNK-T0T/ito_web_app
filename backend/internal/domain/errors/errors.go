package errors

type ErrorCode string

const (
	ValidationError     ErrorCode = "VALIDATION_ERROR"
	AuthenticationError ErrorCode = "AUTHENTICATION_ERROR"
	NotFoundError       ErrorCode = "NOT_FOUND_ERROR"
	BadRequestError     ErrorCode = "BAD_REQUEST_ERROR"
)

type AppError struct {
	Code    ErrorCode
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

// エラー定義
var (
	ErrInvalidInput = AppError{Code: ValidationError, Message: "入力値が不正です"}
	ErrUnauthorized = AppError{Code: AuthenticationError, Message: "認証に失敗しました"}
	ErrNotFound     = AppError{Code: NotFoundError, Message: "リソースが見つかりません"}
)