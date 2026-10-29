package common

import (
	"flag"
	"net"
)

var (
	Name = "review-service"
)

//===== 下面的从 .env 更新 =====

var (
	PortGRPC = *flag.Int("port", 9090, "the listening port")
	Port     = 8080
	IP       = "127.0.0.1"
)
var (
	DebugMode = false
)

var (
	StartTime = "2026-10-01"
	MachineID = 1
)

func initCommonEnv() {
	{
		PortGRPC = GetEnvOrDefault("PORT_GRPC", PortGRPC)
		Port = GetEnvOrDefault("PORT", Port)
		{ // 尝试获得 IP
			ip, err := getOutboundIP()
			if err == nil {
				IP = ip.String()
			}
		}
		IP = GetEnvOrDefaultString("IP", IP)
	}
	{
		DebugMode = GetEnvOrDefaultBool("DEBUG_MODE", DebugMode)
	}

	{
		StartTime = GetEnvOrDefaultString("START_TIME", StartTime)
		MachineID = GetEnvOrDefault("MACHINE_ID", MachineID)
	}
}

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
