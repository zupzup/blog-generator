package generator

import (
	"fmt"
	"github.com/beevik/etree"
	"os"
	"path/filepath"
	"time"
)

// RSSGenerator object
type RSSGenerator struct {
	Config *RSSConfig
}

// RSSConfig holds the configuration for an RSS feed
type RSSConfig struct {
	Posts       []*Post
	Destination string
}

const rssDateFormat string = "02 Jan 2006 15:04 -0700"

// Generate creates an RSS feed
func (g *RSSGenerator) Generate() error {
	fmt.Println("\tGenerating RSS...")
	posts := g.Config.Posts
	destination := g.Config.Destination
	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	rss := doc.CreateElement("rss")
	rss.CreateAttr("xmlns:atom", "http://www.w3.org/2005/Atom")
	rss.CreateAttr("version", "2.0")
	channel := rss.CreateElement("channel")

	channel.CreateElement("title").SetText(blogTitle)
	channel.CreateElement("link").SetText(blogURL)
	channel.CreateElement("language").SetText(blogLanguage)
	channel.CreateElement("description").SetText(blogDescription)
	channel.CreateElement("lastBuildDate").SetText(time.Now().Format(rssDateFormat))

	atomLink := channel.CreateElement("atom:link")
	atomLink.CreateAttr("href", fmt.Sprintf("%s/index.xml", blogURL))
	atomLink.CreateAttr("rel", "self")
	atomLink.CreateAttr("type", "application/rss+xml")

	for _, post := range posts {
		if err := addItem(channel, post, fmt.Sprintf("%s/%s/", blogURL, post.Name[1:])); err != nil {
			return err
		}
	}

	filePath := filepath.Join(destination, "index.xml")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", filePath, err)
	}
	f.Close()
	if err := doc.WriteToFile(filePath); err != nil {
		return fmt.Errorf("error writing to file %s: %v", filePath, err)
	}
	fmt.Println("\tFinished generating RSS...")
	return nil
}

func addItem(element *etree.Element, post *Post, path string) error {
	meta := post.Meta
	item := element.CreateElement("item")
	item.CreateElement("title").SetText(meta.Title)
	item.CreateElement("link").SetText(path)
	item.CreateElement("guid").SetText(path)
	pubDate, err := time.Parse(dateFormat, meta.Date)
	if err != nil {
		return fmt.Errorf("error parsing date %s: %v", meta.Date, err)
	}
	item.CreateElement("pubDate").SetText(pubDate.Format(rssDateFormat))
	item.CreateElement("description").SetText(string(post.HTML))
	return nil
}
