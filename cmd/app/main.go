package main

import (
	"fmt"
	"log"

	"business/config"
	"business/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	tets := "tets"
	abc := "aaa"
	fmt.Printf("%v \n", tets)

	app.Run(cfg)
}
