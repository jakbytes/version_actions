package github

import (
	"context"
	"testing"
	"version_actions/internal/mocks"

	"github.com/google/go-github/v58/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLatestVersion(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{},
		Ctx:                 context.Background(),
	}
	version, err := repository.LatestVersion()
	require.Nil(t, err)
	require.NotNil(t, version)
	require.Equal(t, "v1.0.1", *version.Name)
}

func TestLatestVersion_Error(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{
			Inner: assert.AnError,
		},
		Ctx: context.Background(),
	}
	_, err := repository.LatestVersion()
	require.NotNil(t, err)
	require.Equal(t, assert.AnError, err)
}

func TestLatestVersion_NoTags(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{
			Tags: []*github.RepositoryTag{},
		},
		Ctx: context.Background(),
	}
	_, err := repository.LatestVersion()
	require.NotNil(t, err)
	require.Equal(t, NoReleaseVersionFound{}, err)
}

func TestLatestVersion_InvalidTag(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{
			Tags: []*github.RepositoryTag{
				{
					Name: github.String("invalid"),
				},
			},
		},
		Ctx: context.Background(),
	}
	_, err := repository.LatestVersion()
	require.NotNil(t, err)
	require.Equal(t, NoReleaseVersionFound{}, err)
}

func TestLatestPrereleaseVersion(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{
			Tags: []*github.RepositoryTag{
				{
					Name: github.String("v1.0.1-alpha.1"),
				},
			},
		},
		Ctx: context.Background(),
	}
	version, err := repository.LatestPrereleaseVersion("alpha")
	require.Nil(t, err)
	require.NotNil(t, version)
	require.Equal(t, "v1.0.1-alpha.1", *version.Name)
}

func TestLatestPrereleaseVersion_Error(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{
			Inner: assert.AnError,
		},
		Ctx: context.Background(),
	}
	_, err := repository.LatestPrereleaseVersion("alpha")
	require.NotNil(t, err)
	require.Equal(t, assert.AnError, err)
}

func TestLatestPrereleaseVersion_NoTags(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{},
		Ctx:                 context.Background(),
	}
	_, err := repository.LatestPrereleaseVersion("alpha")
	require.NotNil(t, err)
	require.Equal(t, NoPrereleaseVersionFound{}, err)
}

func TestPreviousVersion(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{},
		Ctx:                 context.Background(),
	}
	version, err := repository.PreviousVersion()
	require.Nil(t, err)
	require.NotNil(t, version)
	require.Equal(t, "v1.0.0", *version.Name)
}

func TestPreviousVersion_Error(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{
			Inner: assert.AnError,
		},
		Ctx: context.Background(),
	}
	_, err := repository.PreviousVersion()
	require.NotNil(t, err)
	require.Equal(t, assert.AnError, err)
}

func TestPreviousVersion_NoTags(t *testing.T) {
	repository := &Repository{
		RepositoriesService: &mocks.RepositoryService{
			Tags: []*github.RepositoryTag{},
		},
		Ctx: context.Background(),
	}
	_, err := repository.PreviousVersion()
	require.NotNil(t, err)
	require.Equal(t, NoReleaseVersionFound{}, err)
}
