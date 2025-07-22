package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/DanStough/fang"

	"github.com/DanStough/epok/internal/buildinfo"
	"github.com/DanStough/epok/internal/cmd"
	"github.com/DanStough/epok/internal/styles"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

// execute adds all child commands to the root command and sets flags appropriately.
func execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()

	rootCmd := cmd.NewRootCMD()
	theme := styles.NewEpokTheme()
	return fang.Execute(ctx, rootCmd,
		fang.WithCommit(buildinfo.GetCommit()),
		fang.WithVersion(buildinfo.GetVersion()),
		fang.WithColorSchemeFunc(theme.FangColorScheme),
	)
}
