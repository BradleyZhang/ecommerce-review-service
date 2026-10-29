package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"review-service/common"
	"review-service/internal/biz"
	"review-service/internal/data"
	"review-service/internal/server"
	"review-service/internal/service"
	"review-service/pkg/snowflake"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	InitResources()

	// DI
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", common.PortGRPC))
	if err != nil {
		common.FatalLog("failed to listen: ", err)
	}
	db, err := data.NewDB(common.DSN)
	if err != nil {
		common.FatalLog("数据库连接失败:", err)
	}
	d, cleanup, err := data.NewData(db, common.SysLogger)
	defer cleanup()
	if err != nil {
		common.FatalLog("Data 初始化失败", err)
	}
	repo := data.NewReviewRepo(d, common.SysLogger)
	uc := biz.NewReviewUsecase(repo)
	srv := service.NewReviewService(uc)
	grpcServer := server.NewGRPCServer(srv)
	{ // gRPC
		go func() {
			err = grpcServer.Serve(lis)
			if err != nil {
				common.FatalLog("gRPC 错误: ", err)
			}
		}()
		common.SysLog("reivew service gRPC started,listening on ", common.PortGRPC)
	}

	{ // HTTP
		go func() {
			mux, err := server.NewHTTPHandler(ctx)
			if err != nil {
				common.FatalLog("new http server error", err)
			}
			err = http.ListenAndServe(":"+strconv.Itoa(common.Port), mux)
			if err != nil {
				common.FatalLog("http server error: ", err)
			}
		}()
		common.SysLog("review service http started ,listening on ", common.Port)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	common.SysLog("received signal: ", sig, " shutting down...")
}
func InitResources() {
	err := godotenv.Load(".env")
	if err != nil {
		common.SysLog("No .env file found,useing default environment variables.If needed, please create a .env file and set the relevent variables.")
		// common.FatalLog(err)
	}
	common.InitEnv()
	err = snowflake.Init(common.StartTime, int64(common.MachineID))
	if err != nil {
		common.FatalLog("snowflake 初始化失败: ", err)
	}
}
