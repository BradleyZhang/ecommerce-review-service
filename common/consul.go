package common

var (
	ConsulAddress = "localhost:8500"
)

func initConsulEnv() {
	ConsulAddress = GetEnvOrDefaultString("CONSUL_ADDRESS", ConsulAddress)
}
