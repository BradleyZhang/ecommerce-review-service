package common

var (
	DSN = "root:root1234@tcp(127.0.0.1:3306)/review?charset=utf8mb4&parseTime=True"
)

func initDBEnv() {
	DSN = GetEnvOrDefaultString("DSN", DSN)
}
