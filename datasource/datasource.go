package datasource

import (
	"github.com/zupzup/blog-generator/config"
)

// DataSource is the data-source fetching interface
type DataSource interface {
	Fetch(cfg *config.Config) ([]string, error)
}

// New creates a new GitDataSource
func New() DataSource {
	return &GitDataSource{}
}
