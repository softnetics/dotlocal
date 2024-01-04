package main

import (
	"log"

	"github.com/softnetics/dotlocal/internal/daemon"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	dotlocal, err := daemon.NewDotLocal(logger)
	if err != nil {
		panic(err)
	}
	err = dotlocal.Start()
	if err != nil {
		panic(err)
	}

	// err = dotlocal.SetMappings([]internal.Mapping{
	// 	{
	// 		Host:       "app1.local",
	// 		PathPrefix: "/",
	// 		Target:     "http://localhost:3000",
	// 	},
	// 	{
	// 		Host:       "app1.local",
	// 		PathPrefix: "/_api",
	// 		Target:     "http://localhost:4000",
	// 	},
	// 	{
	// 		Host:       "app2.local",
	// 		PathPrefix: "/",
	// 		Target:     "http://localhost:5555",
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	err = dotlocal.Wait()
	if err != nil {
		panic(err)
	}
}
