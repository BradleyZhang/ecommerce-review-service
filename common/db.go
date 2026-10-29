package common

var (
	DSN = ""
)

func initDBEnv() {
	DSN = GetEnvOrDefaultString("DSN", DSN)
}
