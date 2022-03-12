package main

import (
	"context"
	"fmt"
	stdlog "log"

	"test.com/grades"
	"test.com/log"
	"test.com/registry"
	"test.com/service"
)

func main() {
	host, port := "localhost", "6000"
	serviceAddress := fmt.Sprintf("http://%s:%s", host, port)
	reg := registry.Registration{
		ServiceName:      registry.GradingService,
		ServiceURL:       serviceAddress,
		RequiredServices: []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: serviceAddress + "/services",
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
	//  找到提供服务的地址
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Loggin service found at: %s\n", logProvider)
		log.SetClientLogger(logProvider, reg.ServiceName)
	}
	<-ctx.Done()
	fmt.Println("Shutting down grading service")
}
