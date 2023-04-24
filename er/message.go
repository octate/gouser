package er

var messages = map[string]string{
	"1": "Oops! Something went wrong. Please try later",

	"422": "Request not valid",
	"423": "User already exists",
}

var codes = map[Code]string{
	UncaughtException: "1",

	InvalidRequestBody: "422",
	UserAlreadyExists:  "423",
}
