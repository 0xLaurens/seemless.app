package data

// UserStoreErr is a wrapper type related to user storage errors
type UserStoreErr string

var UserStoreError = struct {
	DuplicateUsername UserStoreErr
	NotFound          UserStoreErr
}{
	DuplicateUsername: "DuplicateUsername",
	NotFound:          "NotFound",
}

func UserStoreErrMessage(error UserStoreErr) string {
	switch error {
	case UserStoreError.NotFound:
		return "user not found"
	case UserStoreError.DuplicateUsername:
		return "username not unique"
	default:
		return "internal error"
	}
}

// WsErr is a wrapper type for WS related errors
type WsErr int

// WsError NOTE: should be up-to-date with http status codes
var WsError = struct {
	InvalidRequestBody WsErr
}{
	InvalidRequestBody: 400,
}

func WsErrorType(error WsErr) string {
	switch error {
	case WsError.InvalidRequestBody:
		return "InvalidRequest"
	default:
		return "WsError"
	}
}

func WsErrorMessage(error WsErr) string {
	switch error {
	case WsError.InvalidRequestBody:
		return "invalid request body"
	default:
		return "internal error"
	}
}
