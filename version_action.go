package main

import (
	"github.com/rs/zerolog/log"
	"os"
	"version_actions/action/extract_commit"
	"version_actions/action/pull_request"
	"version_actions/action/version"
	"version_actions/internal/logger"
)

func main() {
	log.Logger = logger.Base()
	action := os.Args[1]
	switch action {
	case "release":
		log.Info().Msg("Release action")
	case "version":
		log.Info().Msg("Version action")
		version.Execute()
	case "pull_request":
		log.Info().Msg("Pull request action")
		pull_request.Execute()
	case "extract_commit":
		log.Info().Msg("Extract commit action")
		extract_commit.ExtractCommit()
	}
}
