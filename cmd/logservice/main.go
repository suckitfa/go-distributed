package main

import (
	"context"
	"fmt"
	stdlog "log"

	"test.com/log"
	"test.com/service"
)

// 日志模块的启动入口
func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	ctx, err := service.Start(
		context.Background(),
		"Log Service",
		host,
		port,
		log.RegisterHandlers,
	)
	if err != nil {
		stdlog.Fatalln(err)
	}
	<-ctx.Done()
	fmt.Println("shutting down log service.")
}
