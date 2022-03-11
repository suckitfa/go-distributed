package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"test.com/registry"
)

func main() {
	http.Handle("/services", &registry.RegistryService{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		log.Println("Server started ! Press any key to shut down....")
		// 按任意键关闭服务
		var s string
		fmt.Scan(&s)
		srv.Shutdown(ctx)
		cancel()
	}()
	<-ctx.Done()
	fmt.Println("Shutting down server....")
}
