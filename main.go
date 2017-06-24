package main

import (
	"flag"

	"github.com/zupzup/blog-generator/cli"
	"github.com/zupzup/blog-generator/config"
)

func init() {
	flag.StringVar(&config.RepoURL, "repo", "", "Source Repository URL")
	flag.StringVar(&config.TmpFolder, "tmpfolder", "./tmp", "folder where the data-source repo is checked out to")
	flag.StringVar(&config.DestFolder, "destfolder", "./www", "is the output folder of the static blog")
}

func main() {
	flag.Parse()
	cli.Run()
}
