package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"sort"
	"strings"
)

// Tag holds the data for a Tag
type Tag struct {
	Name  string
	Link  string
	Count int
}

// ByCountDesc sorts the tags
type ByCountDesc []*Tag

// TagsGenerator object
type TagsGenerator struct {
	Config *TagsConfig
}

// TagsConfig holds the tag's config
type TagsConfig struct {
	TagPostsMap map[string][]*Post
	Template    *template.Template
	Destination string
	Writer      *IndexWriter
}

// Generate creates the tags page
func (g *TagsGenerator) Generate() error {
	fmt.Println("\tGenerating Tags...")
	tagPostsMap := g.Config.TagPostsMap
	t := g.Config.Template
	destination := g.Config.Destination
	tagsPath := filepath.Join(destination, "tags")
	if err := clearAndCreateDestination(tagsPath); err != nil {
		return err
	}
	if err := generateTagIndex(tagPostsMap, t, tagsPath, g.Config.Writer); err != nil {
		return err
	}
	for tag, tagPosts := range tagPostsMap {
		tagPagePath := filepath.Join(tagsPath, tag)
		if err := generateTagPage(tag, tagPosts, t, tagPagePath, g.Config.Writer); err != nil {
			return err
		}
	}
	fmt.Println("\tFinished generating Tags...")
	return nil
}

func generateTagIndex(tagPostsMap map[string][]*Post, t *template.Template, destination string, writer *IndexWriter) error {
	tagsTemplatePath := filepath.Join("static", "tags.html")
	tmpl, err := getTemplate(tagsTemplatePath)
	if err != nil {
		return err
	}
	tags := []*Tag{}
	for tag, posts := range tagPostsMap {
		tags = append(tags, &Tag{Name: tag, Link: getTagLink(tag), Count: len(posts)})
	}
	sort.Sort(ByCountDesc(tags))
	buf := bytes.Buffer{}
	if err := tmpl.Execute(&buf, tags); err != nil {
		return fmt.Errorf("error executing template %s: %v", tagsTemplatePath, err)
	}
	if err := writer.WriteIndexHTML(destination, "Tags", "Tags", template.HTML(buf.String()), t, ""); err != nil {
		return err
	}
	return nil
}

func generateTagPage(tag string, posts []*Post, t *template.Template, destination string, writer *IndexWriter) error {
	if err := clearAndCreateDestination(destination); err != nil {
		return err
	}
	lg := ListingGenerator{&ListingConfig{
		Posts:       posts,
		Template:    t,
		Destination: destination,
		PageTitle:   tag,
		Writer:      writer,
	}}
	if err := lg.Generate(); err != nil {
		return err
	}
	return nil
}

func createTags(tags []string) []*Tag {
	var result []*Tag
	for _, tag := range tags {
		result = append(result, &Tag{Name: tag, Link: getTagLink(tag)})
	}
	return result
}

func getTagLink(tag string) string {
	return fmt.Sprintf("/tags/%s/", strings.ToLower(tag))
}

func (t ByCountDesc) Len() int {
	return len(t)
}

func (t ByCountDesc) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByCountDesc) Less(i, j int) bool {
	return t[i].Count > t[j].Count
}
