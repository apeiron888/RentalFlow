package models

type Error struct{
	Message string
	Code string
}

func (e Error) Error() string {
	return e.Message
}

var (
	USER_NOT_FOUND Error = Error{ Code:"001", Message: "User ID doesn't exist"}
	ADDRESS_NOT_FOUND Error = Error{ Code: "002", Message: "Address ID doesn't exist"}
	PROFILE_PICTURE_TOO_LARGE Error = Error{ Code: "003", Message: "File size should be less than 5MB"}
)