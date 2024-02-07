package pull_request

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"version_actions/internal/logger"
	"version_actions/tools/changelog"
	"version_actions/tools/conventional"
	"version_actions/tools/github"
)

var NewClient = github.NewClient

type Args struct {
	Action string
	Token  string
	Owner  string
	Name   string
	Head   string
	Base   string
}

func getArgs() Args {
	args := os.Args[1:]

	if len(args) < 6 {
		panic("Usage: program token owner name head base")
	}

	return Args{
		Action: args[0],
		Token:  args[1],
		Owner:  args[2],
		Name:   args[3],
		Head:   args[4],
		Base:   args[5],
	}
}

// composeTitle composes a pull request title based on the commit message.
func composeTitle(branch *github.Branch) (title string, err error) {
	title, err = branch.GetLastCommitMessage()
	if err != nil {
		return
	}

	if len(title) > 70 {
		title = title[:70] + "..."
	}
	return title, nil
}

func composeBody(head *github.Branch, base string, existing *string) (body changelog.Markdown, err error) {
	commits, err := head.GetDistinctCommits(base)
	if err != nil {
		return
	}
	pc := conventional.ParseCommits(commits)
	cl := changelog.GenerateNewChangelog(head.RepositoryMetadata.Owner, head.RepositoryMetadata.Name, nil, nil, pc, true)
	if existing == nil { // Create a new body
		body = append(changelog.Markdown{
			"### :robot: I have created a pull request *beep* *boop*",
			"",
			"### Notes",
			"",
			"You can add your personal notes here (above the 'Changelog' section). To ensure your notes and the " +
				"automated changelog updates are maintained correctly, keep the 'Changelog' marker in place. If the " +
				"'Changelog' marker is removed, the automated updates to the changelog will not occur. Personal notes " +
				"above the 'Changelog' will be retained during updates, while content below it will be updated with each " +
				"new commit.",
			"",
		},
			cl...)

		body = append(body, "#", "",
			"This Changelog was composed by [version-action](https://github.com/jakbytes/version-action)",
		)

		return
	} else {
		return updateBody(existing, cl), nil
	}
}

func updateBody(body *string, changelog []string) (lines []string) {
	// Process the existing PR body to retain notes above "Changelog"
	if body != nil && strings.Contains(*body, "## Changelog") {
		existing := strings.Split(*body, "## Changelog")[0]
		if len(existing) != 0 && existing[len(existing)-1] == '\n' {
			existing = existing[:len(existing)-1]
		}
		return append(strings.Split(existing, "\n"), changelog...)
	} else {
		return strings.Split(*body, "\n")
	}
}

func setPullRequest() error {
	args := getArgs()
	ctx := context.Background()
	client := NewClient(ctx, args.Token, args.Owner, args.Name)
	repository := client.Repository()
	head, err := repository.Branch(args.Head)
	if err != nil {
		return fmt.Errorf("failed to get branch: %w", err)
	}

	var title string
	pr, err := client.GetPullRequest(args.Head, args.Base)
	if errors.Is(err, github.NoPullRequestFoundError{Head: args.Head, Base: args.Base}) {
		title, err = composeTitle(head)
		if err != nil {
			return fmt.Errorf("failed to compose pull request title: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to get pull request: %w", err)
	} else {
		title = *pr.Title
	}

	return client.SetPullRequest(args.Head, args.Base, title, true, func(body *string) (changelog.Markdown, error) {
		return composeBody(head, args.Base, body)
	})
}

func Execute() {
	log.Logger = logger.Base()
	err := setPullRequest()
	if err != nil {
		panic(err)
	}
}
