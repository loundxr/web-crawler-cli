package main

import (
	"fmt"
	"log"
	"os"

	"github.com/loundxr/web-crawler-cli/internal/config"
)

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatal("parse config error:", err)
	}
	fmt.Print(cfg)
}
