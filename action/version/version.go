package version

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"version_actions/internal/logger"
	"version_actions/tools"
	"version_actions/tools/github"
	"version_actions/tools/github/composite"
)

var NewClient = github.NewClient

type Args struct {
	Action               string
	Token                string
	Owner                string
	Name                 string
	Head                 string
	Base                 string
	PrereleaseIdentifier string
	ReleaseBranch        string
	Trigger              string
	CommitFiles          []string
}

func setup() (client *github.Client, args Args, err error) {
	input := os.Args[1:]

	if len(input) < 9 {
		panic("Usage: program token owner name head base")
	}

	args = Args{
		Action:               input[0],
		Token:                input[1],
		Owner:                input[2],
		Name:                 input[3],
		Head:                 input[4],
		Base:                 input[5],
		PrereleaseIdentifier: input[6],
		ReleaseBranch:        input[7],
		Trigger:              input[8],
		CommitFiles:          input[9:],
	}

	client = NewClient(context.Background(), args.Token, args.Owner, args.Name)

	if args.ReleaseBranch == "." { // . is the default value for the release branch
		var branch *github.Branch
		branch, err = client.Repository().DefaultBranch()
		if err != nil {
			return nil, args, fmt.Errorf("failed to get default branch: %w", err)
		}
		args.ReleaseBranch = branch.Name
	}

	return
}

func version() {
	client, args, err := setup()
	if err != nil {
		panic(err)
	}

	h := &composite.Handler{
		Client:               client,
		Owner:                args.Owner,
		Name:                 args.Name,
		Head:                 args.Head,
		Base:                 args.Base,
		PrereleaseIdentifier: args.PrereleaseIdentifier,
		ReleaseBranch:        args.ReleaseBranch,
		Trigger:              args.Trigger,
		CommitFiles:          args.CommitFiles,
	}
	err = h.PullRequest()
	if err != nil {
		panic(err)
	}

	tools.OpenOutput(func(out tools.Output) {
		log.Debug().Msgf("Setting version to v%s", h.NextVersion().String())
		out.Set("version", github.String("v"+h.NextVersion().String()))
	})
}

func Execute() {
	log.Logger = logger.Base()
	version()
}
