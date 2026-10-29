package common

import (
	"flag"
)

var (
	PortGRPC = *flag.Int("port", 9090, "the listening port")
	Port     = 8080
)
var (
	DebugMode = false
)

var (
	StartTime = "2026-10-01"
	MachineID = 1
)

func initCommonEnv() {
	PortGRPC = GetEnvOrDefault("PORT_GRPC", PortGRPC)
	Port = GetEnvOrDefault("PORT", Port)

	DebugMode = GetEnvOrDefaultBool("DEBUG_MODE", DebugMode)

	StartTime = GetEnvOrDefaultString("START_TIME", StartTime)
	MachineID = GetEnvOrDefault("MACHINE_ID", MachineID)
}
