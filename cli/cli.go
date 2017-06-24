package cli

import (
	"log"

	"github.com/zupzup/blog-generator/config"
	"github.com/zupzup/blog-generator/datasource"
	"github.com/zupzup/blog-generator/generator"
)

// Run runs the application
func Run() {
	ds := datasource.New()
	dirs, err := ds.Fetch(config.RepoURL, config.TmpFolder)

	if err != nil {
		log.Fatal(err)
	}

	g := generator.New(&generator.SiteConfig{
		Sources:     dirs,
		Destination: config.DestFolder,
	})

	err = g.Generate()

	if err != nil {
		log.Fatal(err)
	}
}
