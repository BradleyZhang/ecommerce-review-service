package common

var (
	ESAddresses []string
	ESIndex     = "review"
	ESAPIKey    string
)

func initESEnv() {
	ESAddresses = GetEnvOrDefaultStrings("ES_ADDRESSES", ESAddresses)
	ESAPIKey = GetEnvOrDefaultString("ES_API_KEY", ESAPIKey)
}
