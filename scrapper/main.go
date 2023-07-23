package main

import (
	"github.com/cloudmachinery/movie-reviews/internal/log"

	"github.com/cloudmachinery/movie-reviews/scrapper/cmd"
	"github.com/spf13/cobra"
)

func main() {
	logger, err := log.SetupLogger(true, "debug")
	if err != nil {
		panic(err)
	}

	root := &cobra.Command{
		Use:   "scrapper",
		Short: "Use this tool to scrap movie info",
	}

	root.AddCommand(cmd.NewScrapCmd(logger))
	root.AddCommand(cmd.NewIngestCmd(logger))

	err = root.Execute()
	if err != nil {
		logger.With("err", err).Error("error executing a command")
	}
}
