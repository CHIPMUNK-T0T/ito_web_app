package errors

type WebSocketError struct {
	Code    string
	Message string
}

func (e WebSocketError) Error() string {
	return e.Message
}

var (
	ErrInvalidMessageFormat = WebSocketError{Code: "INVALID_MESSAGE_FORMAT", Message: "無効なメッセージ形式です"}
	ErrMessageHandling      = WebSocketError{Code: "MESSAGE_HANDLING_ERROR", Message: "メッセージの処理中にエラーが発生しました"}
)