package data

type UserStoreErr string

var UserStoreError = struct {
	DuplicateUsername UserStoreErr
	NotFound          UserStoreErr
}{
	DuplicateUsername: "username is not unique",
	NotFound:          "user does not exist",
}
