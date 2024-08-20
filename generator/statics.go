package generator

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// StaticsGenerator object
type StaticsGenerator struct {
	Config *StaticsConfig
}

// StaticsConfig holds the data for the static sites
type StaticsConfig struct {
	FileToDestination map[string]string
	TemplateToFile    map[string]string
	Template          *template.Template
	Writer            *IndexWriter
}

// Generate creates the static pages
func (g *StaticsGenerator) Generate() error {
	fmt.Println("\tCopying Statics...")
	fileToDestination := g.Config.FileToDestination
	templateToFile := g.Config.TemplateToFile
	t := g.Config.Template
	for k, v := range fileToDestination {
		if err := createFolderIfNotExist(getFolder(v)); err != nil {
			return err
		}
		if err := copyFile(k, v); err != nil {
			return err
		}
	}
	for k, v := range templateToFile {
		if err := createFolderIfNotExist(getFolder(v)); err != nil {
			return err
		}
		content, err := ioutil.ReadFile(k)
		if err != nil {
			return fmt.Errorf("error reading file %s: %v", k, err)
		}
		if err := g.Config.Writer.WriteIndexHTML(getFolder(v), getTitle(k), getTitle(k), template.HTML(content), t, "", ""); err != nil {
			return err
		}
	}
	fmt.Println("\tFinished copying statics...")
	return nil
}

func createFolderIfNotExist(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(path, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory %s: %v", path, err)
			}
		} else {
			return fmt.Errorf("error accessing directory %s: %v", path, err)
		}
	}
	return nil
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", src, err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", dst, err)
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("error writing file %s: %v", dst, err)
	}
	if err := out.Sync(); err != nil {
		return fmt.Errorf("error writing file %s: %v", dst, err)
	}
	return nil
}

func getFolder(path string) string {
	return filepath.Dir(path)
}

func getTitle(path string) string {
	ext := filepath.Ext(path)
	name := filepath.Base(path)
	fileName := name[:len(name)-len(ext)]
	return fmt.Sprintf("%s%s", strings.ToUpper(string(fileName[0])), fileName[1:])
}
