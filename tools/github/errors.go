package github

import "fmt"

type NoReleaseVersionFound struct{}

func (e NoReleaseVersionFound) Error() string {
	return "no commits found with a release tag"
}

type NoPrereleaseVersionFound struct{}

func (e NoPrereleaseVersionFound) Error() string {
	return "no commits found with a prerelease tag"
}

type NoPullRequestFoundError struct {
	Head string
	Base string
}

func (e NoPullRequestFoundError) Error() string {
	return fmt.Errorf("no pull request found for branch %s targeting %s", e.Head, e.Base).Error()
}

type MultiplePullRequestsFoundError struct {
	Head string
	Base string
}

func (e MultiplePullRequestsFoundError) Error() string {
	return fmt.Errorf("multiple pull requests found for branch %s targeting %s", e.Head, e.Base).Error()
}

type BranchNotFound struct{
	Name string
}

func (e BranchNotFound) Error() string {
	return fmt.Errorf("branch %s not found", e.Name).Error()
}
