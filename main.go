package main

import (
	"context"

	appctx "github.com/brave-intl/bat-go/libs/context"
	"github.com/brave/go-translate/server"
)

var (
	// variables will be overwritten at build time
	version   string
	commit    string
	buildTime string
)

func main() {
	ctx := context.Background()

	ctx = context.WithValue(ctx, appctx.VersionCTXKey, version)
	ctx = context.WithValue(ctx, appctx.CommitCTXKey, commit)
	ctx = context.WithValue(ctx, appctx.BuildTimeCTXKey, buildTime)

	server.StartServer(ctx)
}
