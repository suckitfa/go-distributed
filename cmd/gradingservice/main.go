package main

import (
	"context"
	"fmt"
	stdlog "log"

	"test.com/grades"
	"test.com/registry"
	"test.com/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	reg := registry.Registration{
		ServiceName: registry.GradingService,
		ServiceURL:  serviceAddress,
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		reg,
		grades.RegisterHandlers,
	)
	if err != nil {
		stdlog.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
