package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
	"github.com/zupzup/blog-generator/config"
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
	HighlightCSS    template.CSS
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
	Config      *config.Config
}

// New creates a new SiteGenerator
func New(config *SiteConfig) *SiteGenerator {
	return &SiteGenerator{Config: config}
}

// Generate starts the static blog generation
func (g *SiteGenerator) Generate() error {
	templatePath := filepath.Join("static", "template.html")
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
		post, err := newPost(path, g.Config.Config.Blog.Dateformat)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}
	sort.Sort(ByDateDesc(posts))
	if err := runTasks(posts, t, destination, g.Config.Config); err != nil {
		return err
	}
	fmt.Println("Finished generating Site...")
	return nil
}

func runTasks(posts []*Post, t *template.Template, destination string, cfg *config.Config) error {
	var wg sync.WaitGroup
	finished := make(chan bool, 1)
	errors := make(chan error, 1)
	pool := make(chan struct{}, 50)
	generators := []Generator{}

	indexWriter := &IndexWriter{
		BlogURL:         cfg.Blog.URL,
		BlogTitle:       cfg.Blog.Title,
		BlogDescription: cfg.Blog.Description,
		BlogAuthor:      cfg.Blog.Author,
	}

	//posts
	for _, post := range posts {
		pg := PostGenerator{&PostConfig{
			Post:        post,
			Destination: destination,
			Template:    t,
			Writer:      indexWriter,
		}}
		generators = append(generators, &pg)
	}
	tagPostsMap := createTagPostsMap(posts)
	// frontpage
	fg := ListingGenerator{&ListingConfig{
		Posts:       posts[:getNumOfPagesOnFrontpage(posts, cfg.Blog.Frontpageposts)],
		Template:    t,
		Destination: destination,
		PageTitle:   "",
		IsIndex:     true,
		Writer:      indexWriter,
	}}
	// archive
	ag := ListingGenerator{&ListingConfig{
		Posts:       posts,
		Template:    t,
		Destination: filepath.Join(destination, "archive"),
		PageTitle:   "Archive",
		IsIndex:     false,
		Writer:      indexWriter,
	}}
	// tags
	tg := TagsGenerator{&TagsConfig{
		TagPostsMap: tagPostsMap,
		Template:    t,
		Destination: destination,
		Writer:      indexWriter,
	}}

	staticURLs := []string{}
	for _, staticURL := range cfg.Blog.Statics.Templates {
		staticURLs = append(staticURLs, staticURL.Dest)
	}
	// sitemap
	sg := SitemapGenerator{&SitemapConfig{
		Posts:       posts,
		TagPostsMap: tagPostsMap,
		Destination: destination,
		BlogURL:     cfg.Blog.URL,
		Statics:     staticURLs,
	}}
	// rss
	rg := RSSGenerator{&RSSConfig{
		Posts:           posts,
		Destination:     destination,
		DateFormat:      cfg.Blog.Dateformat,
		Language:        cfg.Blog.Language,
		BlogURL:         cfg.Blog.URL,
		BlogDescription: cfg.Blog.Description,
		BlogTitle:       cfg.Blog.Title,
	}}
	// statics
	fileToDestination := map[string]string{}
	for _, static := range cfg.Blog.Statics.Files {
		fileToDestination[static.Src] = filepath.Join(destination, static.Dest)
	}
	templateToFile := map[string]string{}
	for _, static := range cfg.Blog.Statics.Templates {
		templateToFile[static.Src] = filepath.Join(destination, static.Dest, "index.html")
	}
	statg := StaticsGenerator{&StaticsConfig{
		FileToDestination: fileToDestination,
		TemplateToFile:    templateToFile,
		Template:          t,
		Writer:            indexWriter,
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

// IndexWriter writer index.html files
type IndexWriter struct {
	BlogTitle       string
	BlogDescription string
	BlogAuthor      string
	BlogURL         string
}

// WriteIndexHTML writes an index.html file
func (i *IndexWriter) WriteIndexHTML(path, pageTitle, metaDescription string, content template.HTML, t *template.Template) error {
	filePath := filepath.Join(path, "index.html")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	defer f.Close()
	metaDesc := metaDescription
	if metaDescription == "" {
		metaDesc = i.BlogDescription
	}
	hlbuf := bytes.Buffer{}
	hlw := bufio.NewWriter(&hlbuf)
	formatter := html.New(html.WithClasses(true))
	formatter.WriteCSS(hlw, styles.MonokaiLight)
	hlw.Flush()
	w := bufio.NewWriter(f)
	td := IndexData{
		Name:            i.BlogAuthor,
		Year:            time.Now().Year(),
		HTMLTitle:       getHTMLTitle(pageTitle, i.BlogTitle),
		PageTitle:       pageTitle,
		Content:         content,
		CanonicalLink:   buildCanonicalLink(path, i.BlogURL),
		MetaDescription: metaDesc,
		HighlightCSS:    template.CSS(hlbuf.String()),
	}
	if err := t.Execute(w, td); err != nil {
		return fmt.Errorf("error executing template %s: %v", filePath, err)
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("error writing file %s: %v", filePath, err)
	}
	return nil
}

func getHTMLTitle(pageTitle, blogTitle string) string {
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
		return nil, fmt.Errorf("error reading template %s: %v", path, err)
	}
	return t, nil
}

func getNumOfPagesOnFrontpage(posts []*Post, numPosts int) int {
	if len(posts) < numPosts {
		return len(posts)
	}
	return numPosts
}

func buildCanonicalLink(path, baseURL string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return fmt.Sprintf("%s/%s/index.html", baseURL, strings.Join(parts[2:], "/"))
	}
	return "/"
}
