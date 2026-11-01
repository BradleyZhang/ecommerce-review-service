package common

var (
	RedisAddress  string = "localhost:6379"
	RedisPassword        = ""
	RedisDB              = 0
)

func initRedisEnv() {
	RedisAddress = GetEnvOrDefaultString("REDIS_ADDRESS", RedisAddress)
	RedisPassword = GetEnvOrDefaultString("REDIS_PASSWORD", RedisPassword)
	RedisDB = GetEnvOrDefault("REDIS_DB", RedisDB)
}
