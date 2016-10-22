package cli

import (
	"github.com/zupzup/blog-generator/datasource"
	"github.com/zupzup/blog-generator/generator"
	"log"
)

// RepoURL is the URL of the data-source repository
const RepoURL string = "git@github.com:zupzup/blog.git"

// TmpFolder is the folder where the data-source repo is checked out to
const TmpFolder string = "./tmp"

// DestFolder is the output folder of the static blog
const DestFolder string = "./www"

// Run runs the application
func Run() {
	ds := datasource.New()
	dirs, err := ds.Fetch(RepoURL, TmpFolder)
	if err != nil {
		log.Fatal(err)
	}
	g := generator.New(&generator.SiteConfig{
		Sources:     dirs,
		Destination: DestFolder,
	})
	err = g.Generate()
	if err != nil {
		log.Fatal(err)
	}
}
