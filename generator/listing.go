package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// ListingData holds the data for the listing page
type ListingData struct {
	Title      string
	Date       string
	Short      string
	Link       string
	TimeToRead string
	Tags       []*Tag
}

// ArchveLinkData holds the data for the archive link template
type ArchveLinkData struct {
	NumPosts int
}

// ListingGenerator Object
type ListingGenerator struct {
	Config *ListingConfig
}

// ListingConfig holds the configuration for the listing page
type ListingConfig struct {
	Posts                  []*Post
	SumAllPosts            int
	Template               *template.Template
	Destination, PageTitle string
	IsIndex                bool
	Writer                 *IndexWriter
}

// Generate starts the listing generation
func (g *ListingGenerator) Generate() error {
	shortTemplatePath := filepath.Join("static", "short.html")
	archiveLinkTemplatePath := filepath.Join("static", "archiveLink.html")
	posts := g.Config.Posts
	t := g.Config.Template
	destination := g.Config.Destination
	pageTitle := g.Config.PageTitle
	short, err := getTemplate(shortTemplatePath)
	if err != nil {
		return err
	}
	var postBlocks []string
	for _, post := range posts {
		meta := post.Meta
		link := fmt.Sprintf("/%s/", post.Name)
		ld := ListingData{
			Title:      meta.Title,
			Date:       meta.Date,
			Short:      meta.Short,
			Link:       link,
			Tags:       createTags(meta.Tags),
			TimeToRead: calculateTimeToRead(string(post.HTML)),
		}
		block := bytes.Buffer{}
		if err := short.Execute(&block, ld); err != nil {
			return fmt.Errorf("error executing template %s: %v", shortTemplatePath, err)
		}
		postBlocks = append(postBlocks, block.String())
	}
	htmlBlocks := template.HTML(strings.Join(postBlocks, "<br />"))
	if g.Config.IsIndex {
		archiveLink, err := getTemplate(archiveLinkTemplatePath)
		if err != nil {
			return err
		}
		lastBlock := bytes.Buffer{}
		ald := ArchveLinkData{
			NumPosts: g.Config.SumAllPosts,
		}
		if err := archiveLink.Execute(&lastBlock, ald); err != nil {
			return fmt.Errorf("error executing template %s: %v", archiveLinkTemplatePath, err)
		}
		htmlBlocks = template.HTML(fmt.Sprintf("%s%s", htmlBlocks, template.HTML(lastBlock.String())))
	} else {
		if strings.Contains(g.Config.Destination, "tags") {
			pageTitle = fmt.Sprintf("%s (%d articles)", pageTitle, len(postBlocks))
		} else {
			pageTitle = fmt.Sprintf("%s (%d articles)", pageTitle, g.Config.SumAllPosts)
		}
	}
	if err := g.Config.Writer.WriteIndexHTML(destination, pageTitle, pageTitle, htmlBlocks, t, "", ""); err != nil {
		return err
	}
	return nil
}

func calculateTimeToRead(input string) string {
	// an average human reads about 200 wpm, but we use a bit more, since a lot of it will be code
	var secondsPerWord = 60.0 / 250.0
	// multiply with the amount of words
	words := secondsPerWord * float64(len(strings.Split(input, " ")))
	// add 5 seconds for each image
	images := 5.0 * strings.Count(input, "<img")
	result := (words + float64(images)) / 60.0
	if result < 1.0 {
		result = 1.0
	}
	return fmt.Sprintf("%.0fm", result)
}
