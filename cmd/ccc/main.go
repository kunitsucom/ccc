package main

import (
	"context"
	"os"

	"github.com/kunitsuinc/ccc/pkg/entrypoint"
	"github.com/kunitsuinc/ccc/pkg/log"
)

func main() {
	ctx := context.Background()

	if err := entrypoint.CCC(ctx); err != nil {
		log.Errorf("entrypoint.CCC: %v", err)
		os.Exit(1)
	}
}
