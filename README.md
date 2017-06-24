[![Go Report Card](https://goreportcard.com/badge/github.com/zupzup/calories)](https://goreportcard.com/report/github.com/zupzup/calories)

# blog-generator

A static blog generator using a configurable GitHub repository as a data-source. The posts are written in markdown with yml metadata attached to them. [This](https://github.com/zupzup/blog) is an example repo for the blog at [https://zupzup.org/](https://zupzup.org/).

## Features

* Listing
* Sitemap Generator
* RSS Feed
* Code Highlighting
* Archive 
* Configurable Static Pages 
* Tags 

## Installation

```bash
go get github.com/zupzup/blog-generator
```

## Usage

Just execute

```bash
blog-generator -repo git@github.com:username/repo.git
```

in this repository's root directory.

## Customization

### Configure the CLI

Provide `flag` to binary. To know which flags are used and which are used and
to override them:

```shell
go build

./blog-generator -h
```

Output:

```
Usage of ./blog-generator:
  -destfolder string
        is the output folder of the static blog (default "./www")
  -repo string
        Repository URL
  -tmpfolder string
        folder where the data-source repo is checked out to (default "./tmp")
```

### Configure the Generator

Set the following constants in `generator/generator.go`
```go
// url of the blog
const blogURL = "https://www.someblog.com"
// blog language
const blogLanguage = "en-us"
// blog description
const blogDescription = "some description..."
// date format
const dateFormat string = "02.01.2006"
// main Template
const templatePath string = "static/template.html"
// blog title
const blogTitle string = "my Blog's Title"
// displayed posts on landing page
const numPostsFrontPage int = 10
```

You can define and configure the different generators in the `runTasks` function within `generator/generator.go`.

### Templates

Edit templates in `static` folder to your needs.

## Example Blog Repository

[Blog](https://github.com/zupzup/blog)
