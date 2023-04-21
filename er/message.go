package er

var messages = map[string]string{
	"1": "Oops! Something went wrong. Please try later",

	"100": "Invalid Application ID or token",
}

var codes = map[Code]string{
	UncaughtException: "1",

	InvalidAppToken: "100",
	UserNotPresent:  "201",
}
