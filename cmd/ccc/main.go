package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kunitsuinc/ccc/pkg/config"
	"github.com/kunitsuinc/ccc/pkg/entrypoint"
	"github.com/kunitsuinc/ccc/pkg/log"
)

func main() {
	ctx := context.Background()
	config.Load()

	if config.SubcommandVersion() {
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n%s\n", config.Version(), config.Revision(), config.Branch(), config.Timestamp())
		return
	}

	config.MustCheck()

	if err := entrypoint.CCC(ctx); err != nil {
		log.Errorf("entrypoint.CCC: %v", err)
		os.Exit(1)
	}
}
