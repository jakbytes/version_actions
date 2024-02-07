package github

import (
	"context"
	"errors"
	"github.com/jakbytes/version_actions/internal/mocks"
	"testing"

	"github.com/google/go-github/v58/github"
	"github.com/stretchr/testify/require"
)

func TestGetBranch(t *testing.T) {
	repository := &Repository{
		branches:            make(map[string]*Branch),
		RepositoriesService: &mocks.RepositoryService{},
		Ctx:                 context.Background(),
	}
	branch, err := repository.getBranch("main")
	require.Nil(t, err)
	require.NotNil(t, branch)
	require.Equal(t, "main", *branch.Name)
}

func TestGetBranch_Error(t *testing.T) {
	repository := &Repository{
		branches: make(map[string]*Branch),
		RepositoriesService: &mocks.RepositoryService{
			GetBranchError: func(ctx context.Context, owner string, repo string, branch string, maxRedirects int) error {
				return errors.New("unexpected status code: 404 Not Found")
			},
		},
		Ctx: context.Background(),
	}
	_, err := repository.getBranch("main")
	require.NotNil(t, err)
	require.Equal(t, BranchNotFound{Name: "main"}, err)
}

func TestBranch_Error(t *testing.T) {
	repository := &Repository{
		branches: make(map[string]*Branch),
		RepositoriesService: &mocks.RepositoryService{
			GetBranchError: func(ctx context.Context, owner string, repo string, branch string, maxRedirects int) error {
				return errors.New("404")
			},
		},
		Ctx: context.Background(),
	}
	_, err := repository.Branch("main")
	require.NotNil(t, err)
	require.Equal(t, BranchNotFound{Name: "main"}, err)
}

func TestDefaultBranch(t *testing.T) {
	repository := &Repository{
		branches:            make(map[string]*Branch),
		RepositoriesService: &mocks.RepositoryService{},
		Ctx:                 context.Background(),
	}
	branch, err := repository.DefaultBranch()
	require.Nil(t, err)
	require.NotNil(t, branch)
	require.Equal(t, "main", branch.Name)
}

func TestDefaultBranch_Error(t *testing.T) {
	repository := &Repository{
		branches: make(map[string]*Branch),
		RepositoriesService: &mocks.RepositoryService{
			GetError: errors.New("404"),
		},
		Ctx: context.Background(),
	}
	branch, err := repository.DefaultBranch()
	require.NotNil(t, err)
	require.Nil(t, branch)
	require.Equal(t, errors.New("404"), err)
}

func TestCreateBranch(t *testing.T) {
	repository := &Repository{
		branches:            make(map[string]*Branch),
		RepositoriesService: &mocks.RepositoryService{},
		GitService:          &mocks.GitService{},
		Ctx:                 context.Background(),
	}
	branch, err := repository.CreateBranch("main", github.String("hash1-hash1"))
	require.Nil(t, err)
	require.NotNil(t, branch)
	require.Equal(t, "main", branch.Name)
}

func TestCreateBranch_Error(t *testing.T) {
	repository := &Repository{
		branches:            make(map[string]*Branch),
		RepositoriesService: &mocks.RepositoryService{},
		GitService: &mocks.GitService{
			CreateRefError: errors.New("404"),
		},
		Ctx: context.Background(),
	}
	_, err := repository.CreateBranch("main", github.String("hash1-hash1"))
	require.NotNil(t, err)
	require.Equal(t, errors.New("404"), err)
}
