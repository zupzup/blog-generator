package cli

import (
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	"github.com/zupzup/blog-generator/config"
	"github.com/zupzup/blog-generator/datasource"
	"github.com/zupzup/blog-generator/generator"
)

// Run runs the application
func Run() {
	cfg, err := readConfig()
	if err != nil {
		log.Fatal("There was an error while reading the configuration file: ", err)
	}
	ds := datasource.New()
	dirs, err := ds.Fetch(cfg.Generator.Repo, cfg.Generator.Tmp)

	if err != nil {
		log.Fatal(err)
	}

	g := generator.New(&generator.SiteConfig{
		Sources:     dirs,
		Destination: cfg.Generator.Dest,
		Config:      cfg,
	})

	err = g.Generate()

	if err != nil {
		log.Fatal(err)
	}
}

func readConfig() (*config.Config, error) {
	data, err := ioutil.ReadFile("bloggen.yml")
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %v", err)
	}
	cfg := config.Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("could not parse config: %v", err)
	}
	if cfg.Generator.Repo == "" {
		return nil, fmt.Errorf("Please provide a repository URL, e.g.: https://github.com/zupzup/blog")
	}
	if cfg.Generator.Tmp == "" {
		cfg.Generator.Tmp = "tmp"
	}
	if cfg.Generator.Dest == "" {
		cfg.Generator.Dest = "www"
	}
	if cfg.Blog.URL == "" {
		return nil, fmt.Errorf("Please provide a Blog URL, e.g.: https://www.zupzup.org")
	}
	if cfg.Blog.Language == "" {
		cfg.Blog.Language = "en-us"
	}
	if cfg.Blog.Description == "" {
		return nil, fmt.Errorf("Please provide a Blog Description, e.g.: A blog about Go, JavaScript, Open Source and Programming in General")
	}
	if cfg.Blog.Dateformat == "" {
		cfg.Blog.Dateformat = "02.01.2006"
	}
	if cfg.Blog.Title == "" {
		return nil, fmt.Errorf("Please provide a Blog Title, e.g.: zupzup")
	}
	if cfg.Blog.Author == "" {
		return nil, fmt.Errorf("Please provide a Blog author, e.g.: Mario Zupan")
	}
	if cfg.Blog.Frontpageposts == 0 {
		cfg.Blog.Frontpageposts = 10
	}
	return &cfg, nil
}
