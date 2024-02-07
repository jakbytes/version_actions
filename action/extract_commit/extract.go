package extract_commit

import (
	"context"
	"fmt"
	"github.com/jakbytes/version_actions/tools"
	"github.com/jakbytes/version_actions/tools/conventional"
	"github.com/jakbytes/version_actions/tools/github"
	"github.com/leodido/go-conventionalcommits"
	cparser "github.com/leodido/go-conventionalcommits/parser"
	"os"
	"strings"
)

var MaxDepth = 10

func ExtractCommit() {
	input := os.Args[2:]
	token := input[0]
	owner := input[1]
	name := input[2]
	branchName := input[3]

	client := github.NewClient(context.Background(), token, owner, name)

	parser := conventional.Parser{Machine: cparser.NewMachine(
		conventionalcommits.WithTypes(conventionalcommits.TypesFreeForm),
	)}

	branch, err := client.Repository().Branch(branchName)
	if err != nil {
		panic(err)
	}

	commits := branch.Commits(1)
	if err != nil {
		panic(err)
	}

	for commits.Next() {
		for _, rawCommit := range commits.Commits() {
			commit := parser.ParseCommit(rawCommit)
			if commit != nil && commit.Ok() {
				tools.OpenOutput(func(out tools.Output) {
					out.Set("type", &commit.Type)
					out.Set("description", &commit.Description)
					out.Set("scope", commit.Scope)
					out.Set("exclamation", tools.String(fmt.Sprintf("%t", commit.Exclamation)))
					out.Set("body", commit.Body)

					// Handle footers; since footers are a map, you might want to join them as a single string
					if len(commit.Footers) > 0 {
						var footers []string
						for key, values := range commit.Footers {
							// Join multiple values for the same footer key with a comma or another separator
							footerValue := strings.Join(values, ", ")
							footers = append(footers, fmt.Sprintf("%s=%s", key, footerValue))
						}
						out.Set("footers", tools.String(strings.Join(footers, "; ")))
					}
				})
				if err != nil {
					panic(err)
				}
				os.Exit(0)
			}
		}
	}
}
