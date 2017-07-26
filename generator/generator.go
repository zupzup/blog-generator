package generator

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Meta is a data container for Metadata
type Meta struct {
	Title      string
	Short      string
	Date       string
	Tags       []string
	ParsedDate time.Time
}

// IndexData is a data container for the landing page
type IndexData struct {
	HTMLTitle       string
	PageTitle       string
	Content         template.HTML
	Year            int
	Name            string
	CanonicalLink   string
	MetaDescription string
}

// Generator interface
type Generator interface {
	Generate() error
}

// SiteGenerator object
type SiteGenerator struct {
	Config *SiteConfig
}

// SiteConfig holds the sources and destination folder
type SiteConfig struct {
	Sources     []string
	Destination string
}

// New creates a new SiteGenerator
func New(config *SiteConfig) *SiteGenerator {
	return &SiteGenerator{Config: config}
}

const blogURL = "https://www.zupzup.org"
const blogLanguage = "en-us"
const blogDescription = "A blog about Go, JavaScript and Programming in General"
const defaultMeta = "A blog about Go, JavaScript, Open Source and Programming in General"
const dateFormat string = "02.01.2006"
const templatePath string = "static/template.html"
const blogTitle string = "zupzup"
const numPostsFrontPage int = 10

// Generate starts the static blog generation
func (g *SiteGenerator) Generate() error {
	fmt.Println("Generating Site...")
	sources := g.Config.Sources
	destination := g.Config.Destination
	if err := clearAndCreateDestination(destination); err != nil {
		return err
	}
	if err := clearAndCreateDestination(filepath.Join(destination, "archive")); err != nil {
		return err
	}
	t, err := getTemplate(templatePath)
	if err != nil {
		return err
	}
	var posts []*Post
	for _, path := range sources {
		post, err := newPost(path)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}
	sort.Sort(ByDateDesc(posts))
	if err := runTasks(posts, t, destination); err != nil {
		return err
	}
	fmt.Println("Finished generating Site...")
	return nil
}

func runTasks(posts []*Post, t *template.Template, destination string) error {
	var wg sync.WaitGroup
	finished := make(chan bool, 1)
	errors := make(chan error, 1)
	pool := make(chan struct{}, 50)
	generators := []Generator{}

	//posts
	for _, post := range posts {
		pg := PostGenerator{&PostConfig{
			Post:        post,
			Destination: destination,
			Template:    t,
		}}
		generators = append(generators, &pg)
	}
	tagPostsMap := createTagPostsMap(posts)
	// frontpage
	fg := ListingGenerator{&ListingConfig{
		Posts:       posts[:getNumOfPagesOnFrontpage(posts)],
		Template:    t,
		Destination: destination,
		PageTitle:   "",
		IsIndex:     true,
	}}
	// archive
	ag := ListingGenerator{&ListingConfig{
		Posts:       posts,
		Template:    t,
		Destination: filepath.Join(destination, "archive"),
		PageTitle:   "Archive",
		IsIndex:     false,
	}}
	// tags
	tg := TagsGenerator{&TagsConfig{
		TagPostsMap: tagPostsMap,
		Template:    t,
		Destination: destination,
	}}

	// sitemap
	sg := SitemapGenerator{&SitemapConfig{
		Posts:       posts,
		TagPostsMap: tagPostsMap,
		Destination: destination,
	}}
	// rss
	rg := RSSGenerator{&RSSConfig{
		Posts:       posts,
		Destination: destination,
	}}
	// statics
	fileToDestination := map[string]string{
		"static/favicon.ico": filepath.Join(destination, "favicon.ico"),
		"static/robots.txt":  filepath.Join(destination, "robots.txt"),
		"static/about.png":   filepath.Join(destination, "about.png"),
	}
	templateToFile := map[string]string{
		"static/about.html": filepath.Join(destination, "/about/index.html"),
	}
	statg := StaticsGenerator{&StaticsConfig{
		FileToDestination: fileToDestination,
		TemplateToFile:    templateToFile,
		Template:          t,
	}}
	generators = append(generators, &fg, &ag, &tg, &sg, &rg, &statg)

	for _, generator := range generators {
		wg.Add(1)
		go func(g Generator) {
			defer wg.Done()
			pool <- struct{}{}
			defer func() { <-pool }()
			if err := g.Generate(); err != nil {
				errors <- err
			}
		}(generator)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
		return nil
	case err := <-errors:
		if err != nil {
			return err
		}
	}
	return nil
}

func clearAndCreateDestination(path string) error {
	if err := os.RemoveAll(path); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("error removing folder at destination %s: %v ", path, err)
		}
	}
	return os.Mkdir(path, os.ModePerm)
}

func writeIndexHTML(path, pageTitle string, metaDescription string, content template.HTML, t *template.Template) error {
	filePath := filepath.Join(path, "index.html")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer f.Close()
	metaDesc := metaDescription
	if metaDescription == "" {
		metaDesc = defaultMeta
	}
	w := bufio.NewWriter(f)
	td := IndexData{
		Name:            "Mario Zupan",
		Year:            time.Now().Year(),
		HTMLTitle:       getHTMLTitle(pageTitle),
		PageTitle:       pageTitle,
		Content:         content,
		CanonicalLink:   buildCanonicalLink(path, blogURL),
		MetaDescription: metaDesc,
	}
	if err := t.Execute(w, td); err != nil {
		return fmt.Errorf("error executing template %s: %v", templatePath, err)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}
	return nil
}

func getHTMLTitle(pageTitle string) string {
	if pageTitle == "" {
		return blogTitle
	}
	return fmt.Sprintf("%s - %s", pageTitle, blogTitle)
}

func createTagPostsMap(posts []*Post) map[string][]*Post {
	result := make(map[string][]*Post)
	for _, post := range posts {
		for _, tag := range post.Meta.Tags {
			key := strings.ToLower(tag)
			if result[key] == nil {
				result[key] = []*Post{post}
			} else {
				result[key] = append(result[key], post)
			}
		}
	}
	return result
}

func getTemplate(path string) (*template.Template, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("error reading template %s: %v", templatePath, err)
	}
	return t, nil
}

func getNumOfPagesOnFrontpage(posts []*Post) int {
	if len(posts) < numPostsFrontPage {
		return len(posts)
	}
	return numPostsFrontPage
}

func buildCanonicalLink(path, baseURL string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return fmt.Sprintf("%s/%s/index.html", baseURL, strings.Join(parts[2:], "/"))
	}
	return "/"
}
