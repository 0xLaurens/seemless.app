package data

// UserStoreErr is a wrapper type related to user storage errors
type UserStoreErr string

var UserStoreError = struct {
	DuplicateUsername UserStoreErr
	NotFound          UserStoreErr
}{
	DuplicateUsername: "username is not unique",
	NotFound:          "user does not exist",
}

// WsErr is a wrapper type for WS related errors
type WsErr string

var WsError = struct {
	InvalidRequestBody WsErr
}{
	InvalidRequestBody: "request body was not formatted properly",
}
