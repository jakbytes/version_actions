package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v58/github"
	"github.com/rs/zerolog/log"
)

// Branch is a struct that contains the RepositoriesService, context, token, owner, name, and branch. It is used to
// interact with the GitHub API in the context of a specific repository.
type Branch struct {
	RepositoriesService
	RepositoryMetadata RepositoryMetadata
	GitService
	*github.Branch
	Ctx  context.Context
	Name string
}

// GetDistinctCommits returns a map of unique commits from the base branch to the head branch
//
// Parameters:
//   - base: The base branch to compare against
//
// Returns:
//   - map[string]*github.RepositoryCommit: A map of unique commits between the base and current branch
//   - error: An error if one occurred
func (b *Branch) GetDistinctCommits(base string) (commits map[string]*github.RepositoryCommit, err error) {
	comparison, _, err := b.CompareCommits(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, base, b.Name, nil)
	if err != nil {
		return nil, err
	}

	commits = make(map[string]*github.RepositoryCommit)

	// Extract unique commits
	for _, commit := range comparison.Commits {
		hash := commit.GetSHA()
		c := &github.RepositoryCommit{}
		*c = *commit
		commits[hash] = c
	}

	return commits, nil
}

// GetCommitsSinceCommit returns a map of commits from the hash commit to the latest commit, if hash is nil, it will
// return all commits on the current branch
func (b *Branch) GetCommitsSinceCommit(hash *string) (map[string]*github.RepositoryCommit, error) {
	commits := make(map[string]*github.RepositoryCommit)
	nextPage := 0
	for {
		pages, response, err := b.ListCommits(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, &github.CommitsListOptions{SHA: b.Name, ListOptions: github.ListOptions{Page: nextPage, PerPage: 10}})
		if err != nil {
			return nil, err
		}
		for _, commit := range pages {
			if hash != nil {
				log.Debug().Msgf("> Commit: %s", *commit.SHA)
				log.Debug().Msgf(">   Hash: %s", *hash)
				log.Debug().Msgf(">  Equal: %t", *commit.SHA == *hash)
			}
			if hash != nil && *commit.SHA == *hash {
				return commits, nil // return if we've reached the hash commit
			}
			commits[*commit.SHA] = commit
		}

		if response != nil {
			if nextPage = response.NextPage; nextPage == 0 {
				break // break if there are no more pages
			}
		}
	}
	return commits, nil
}

// GetLastCommitMessage retrieves the last commit message from the branch
func (b *Branch) GetLastCommitMessage() (string, error) {
	commits, _, err := b.ListCommits(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, &github.CommitsListOptions{
		SHA:         b.Name, // can be any branch or commit SHA
		ListOptions: github.ListOptions{PerPage: 1},
	})
	if err != nil {
		return "", err
	}

	if len(commits) == 0 {
		return "", fmt.Errorf("no commits found in branch")
	}

	return *commits[0].Commit.Message, nil
}

func (b *Branch) Reset(sha *string) error {
	ref := &github.Reference{Ref: github.String("refs/heads/" + b.Name), Object: &github.GitObject{SHA: sha}}
	_, _, err := b.UpdateRef(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, ref, true)
	return err
}

type File struct {
	Path    string
	Content string
}

// AddFiles adds multiple files to a branch and returns the new tree SHA and parent commit SHA
func (b *Branch) AddFiles(files []File) (newTreeSHA string, parentCommitSHA string, err error) {
	var entries []*github.TreeEntry

	// Create a Blob for each file content
	for _, file := range files {
		blob, _, err := b.CreateBlob(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, &github.Blob{
			Content:  github.String(file.Content),
			Encoding: github.String("utf-8"),
		})
		if err != nil {
			return "", "", err
		}

		entries = append(entries, &github.TreeEntry{
			Path: github.String(file.Path),
			Type: github.String("blob"),
			Mode: github.String("100644"),
			SHA:  blob.SHA,
		})
	}

	// Get the latest commit to find the parent commit SHA and the current tree SHA
	commits, _, err := b.ListCommits(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, &github.CommitsListOptions{
		SHA:         b.Name, // Branch name
		ListOptions: github.ListOptions{PerPage: 1},
	})
	if err != nil {
		return "", "", err
	}
	if len(commits) == 0 {
		return "", "", fmt.Errorf("no commits found in branch %s", b.Name)
	}
	latestCommit := commits[0]

	// Create a Tree with the new files
	tree, _, err := b.CreateTree(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, *latestCommit.Commit.Tree.SHA, entries)
	if err != nil {
		return "", "", err
	}

	return *tree.SHA, *latestCommit.SHA, nil
}

func (b *Branch) CommitChanges(newTreeSHA, parentCommitSHA, commitMessage string) error {
	// Create a Commit with the new tree and the parent commit SHA
	commit, _, err := b.CreateCommit(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, &github.Commit{
		Message: github.String(commitMessage),
		Tree:    &github.Tree{SHA: github.String(newTreeSHA)},
		Parents: []*github.Commit{{SHA: github.String(parentCommitSHA)}},
	}, nil)
	if err != nil {
		return err
	}

	// Update the Reference to point to the new commit SHA
	ref := &github.Reference{Ref: github.String("refs/heads/" + b.Name), Object: &github.GitObject{SHA: commit.SHA}}
	_, _, err = b.UpdateRef(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, ref, true)
	return err
}

// CommitIterator iterates over commits in a branch.
type CommitIterator struct {
	*Branch
	perPage int
	page    int
	commits []*github.RepositoryCommit
	err     error
}

func (b *Branch) Commits(perPage int) CommitIterator {
	return CommitIterator{
		Branch:  b,
		perPage: perPage,
		page:    1,
	}
}

// Next loads the next batch of commits if needed and advances the iterator.
func (b *CommitIterator) Next() bool {
	b.fetchCommits()
	return b.page != 0
}

// fetchCommits retrieves commits from GitHub.
func (b *CommitIterator) fetchCommits() {
	commits, response, err := b.ListCommits(b.Ctx, b.RepositoryMetadata.Owner, b.RepositoryMetadata.Name, &github.CommitsListOptions{
		SHA:         b.Name,
		ListOptions: github.ListOptions{PerPage: b.perPage, Page: b.page},
	})
	if err != nil {
		b.err = err
		b.page = 0 // Stop iteration on error
		return
	}
	b.commits = commits
	b.page = response.NextPage
}

// Commits returns the current batch of commits.
func (b *CommitIterator) Commits() []*github.RepositoryCommit {
	return b.commits
}

// Err returns the last error encountered by the iterator.
func (b *CommitIterator) Err() error {
	return b.err
}
