package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"test.com/registry"
)

func Start(
	ctx context.Context,
	host, port string,
	reg registry.Registration,
	registerHandlersFunc func()) (context.Context, error) {
	registerHandlersFunc()
	// 启动服务
	ctx = startService(ctx, reg.ServiceName, host, port)
	// 注册服务
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startService(
	ctx context.Context,
	serviceName registry.ServiceName,
	host, port string) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	var server http.Server
	server.Addr = ":" + port

	// 服务启动出错
	go func() {
		log.Println(server.ListenAndServe())
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		cancel()
	}()

	// 手动停止服务
	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", serviceName)
		var s string
		fmt.Scanln(&s)
		err := registry.ShutdownService(fmt.Sprintf("http://%s:%s", host, port))
		if err != nil {
			log.Println(err)
		}
		server.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
