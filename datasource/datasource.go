package datasource

// DataSource is the data-source fetching interface
type DataSource interface {
	Fetch(from, to string) ([]string, error)
}

// New creates a new GitDataSource
func New() DataSource {
	return &GitDataSource{}
}
