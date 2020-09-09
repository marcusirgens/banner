package main

import (
	"context"
	"github.com/marcusirgens/banner/internal/config"
	"github.com/marcusirgens/banner/printer"
)

func main() {
	c := config.DefaultConfig
	ctx := context.Background()

	if c.GithubAccessToken != "" {
		printer.PrintEvents(ctx, c.GithubAccessToken)
	}
}
