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
type WsErr int

// WsError NOTE: should be up-to-date with http status codes
var WsError = struct {
	InvalidRequestBody WsErr
}{
	InvalidRequestBody: 400,
}
