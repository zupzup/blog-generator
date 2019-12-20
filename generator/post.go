package generator

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Post holds data for a post
type Post struct {
	Name      string
	HTML      []byte
	Meta      *Meta
	ImagesDir string
	Images    []string
}

// ByDateDesc is the sorting object for posts
type ByDateDesc []*Post

// PostGenerator object
type PostGenerator struct {
	Config *PostConfig
}

// PostConfig holds the post's configuration
type PostConfig struct {
	Post        *Post
	Destination string
	Template    *template.Template
	Writer      *IndexWriter
}

// Generate generates a post
func (g *PostGenerator) Generate() error {
	post := g.Config.Post
	destination := g.Config.Destination
	t := g.Config.Template
	fmt.Printf("\tGenerating Post: %s...\n", post.Meta.Title)
	staticPath := filepath.Join(destination, post.Name)
	if err := os.Mkdir(staticPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory at %s: %v", staticPath, err)
	}
	if post.ImagesDir != "" {
		if err := copyImagesDir(post.ImagesDir, staticPath); err != nil {
			return err
		}
	}

	if err := g.Config.Writer.WriteIndexHTML(staticPath, post.Meta.Title, post.Meta.Short, template.HTML(string(post.HTML)), t); err != nil {
		return err
	}
	fmt.Printf("\tFinished generating Post: %s...\n", post.Meta.Title)
	return nil
}

func newPost(path, dateFormat string) (*Post, error) {
	meta, err := getMeta(path, dateFormat)
	if err != nil {
		return nil, err
	}
	html, err := getHTML(path)
	if err != nil {
		return nil, err
	}
	imagesDir, images, err := getImages(path)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(path)

	return &Post{Name: name, Meta: meta, HTML: html, ImagesDir: imagesDir, Images: images}, nil
}

func copyImagesDir(source, destination string) (err error) {
	path := filepath.Join(destination, "images")
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return fmt.Errorf("error creating images directory at %s: %v", path, err)
	}
	files, err := ioutil.ReadDir(source)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", path, err)
	}
	for _, file := range files {
		src := filepath.Join(source, file.Name())
		dst := filepath.Join(path, file.Name())
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func getMeta(path, dateFormat string) (*Meta, error) {
	filePath := filepath.Join(path, "meta.yml")
	metaraw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	meta := Meta{}
	err = yaml.Unmarshal(metaraw, &meta)
	if err != nil {
		return nil, fmt.Errorf("error reading yml in %s: %v", filePath, err)
	}
	parsedDate, err := time.Parse(dateFormat, meta.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date in %s: %v", filePath, err)
	}
	meta.ParsedDate = parsedDate
	return &meta, nil
}

func getHTML(path string) ([]byte, error) {
	filePath := filepath.Join(path, "post.md")
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	html := blackfriday.MarkdownCommon(input)
	replaced, err := replaceCodeParts(html)
	if err != nil {
		return nil, fmt.Errorf("error during syntax highlighting of %s: %v", filePath, err)
	}
	return []byte(replaced), nil

}

func getImages(path string) (string, []string, error) {
	dirPath := filepath.Join(path, "images")
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, nil
		}
		return "", nil, fmt.Errorf("error while reading folder %s: %v", dirPath, err)
	}
	images := []string{}
	for _, file := range files {
		images = append(images, file.Name())
	}
	return dirPath, images, nil
}

func replaceCodeParts(htmlFile []byte) (string, error) {
	byteReader := bytes.NewReader(htmlFile)
	doc, err := goquery.NewDocumentFromReader(byteReader)
	if err != nil {
		return "", fmt.Errorf("error while parsing html: %v", err)
	}
	// find code-parts via css selector and replace them with highlighted versions
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		lang := strings.TrimPrefix(class, "language-")
		oldCode := s.Text()
		lexer := lexers.Get(lang)
		formatter := html.New(html.WithClasses(true))
		iterator, err := lexer.Tokenise(nil, string(oldCode))
		if err != nil {
			fmt.Printf("ERROR during syntax highlighting, %v", err)
		}
		b := bytes.Buffer{}
		buf := bufio.NewWriter(&b)
		err = formatter.Format(buf, styles.GitHub, iterator)
		if err != nil {
			fmt.Printf("ERROR during syntax highlighting, %v", err)
		}
		buf.Flush()
		s.SetHtml(b.String())
	})
	new, err := doc.Html()
	if err != nil {
		return "", fmt.Errorf("error while generating html: %v", err)
	}
	// replace unnecessarily added html tags
	new = strings.Replace(new, "<html><head></head><body>", "", 1)
	new = strings.Replace(new, "</body></html>", "", 1)
	return new, nil
}

func (p ByDateDesc) Len() int {
	return len(p)
}

func (p ByDateDesc) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByDateDesc) Less(i, j int) bool {
	return p[i].Meta.ParsedDate.After(p[j].Meta.ParsedDate)
}
