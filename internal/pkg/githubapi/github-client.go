package githubapi

import (
	"context"

	"github.com/google/go-github/v49/github"
)

var (
	ctx                 = context.Background()
	standardListOptions = github.ListOptions{
		PerPage: 100,
		Page:    1,
	}
)
