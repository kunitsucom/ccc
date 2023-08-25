package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kunitsucom/ccc/pkg/config"
	"github.com/kunitsucom/ccc/pkg/entrypoint"
	"github.com/kunitsucom/ccc/pkg/log"
	"github.com/kunitsucom/util.go/must"
)

func main() {
	ctx := context.Background()
	config.Load()

	if config.SubcommandVersion() {
		fmt.Fprintf(os.Stdout, "%s\n%s\n%s\n%s\n", config.Version(), config.Revision(), config.Branch(), config.Timestamp())
		return
	}

	must.Must(config.Check())

	if err := entrypoint.CCC(ctx); err != nil {
		log.Errorf("entrypoint.CCC: %v", err)
		os.Exit(1)
	}
}
