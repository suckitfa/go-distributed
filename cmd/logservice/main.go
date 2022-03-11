package main

import (
	"context"
	"fmt"
	stdlog "log"

	"test.com/log"
	"test.com/registry"
	"test.com/service"
)

// 日志模块的启动入口
func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	r := registry.Registration{
		ServiceName: "Log Service",
		ServiceURL:  serviceAddress,
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		log.RegisterHandlers,
	)
	if err != nil {
		stdlog.Fatalln(err)
	}

	// 接受返回的contextDone — 返回一个 Channel，
	// 这个 Channel 会在当前工作完成或者上下文被取消后关闭，
	// 多次调用 Done 方法会返回同一个 Channel；
	<-ctx.Done()
	fmt.Println("shutting down log service.")
}
