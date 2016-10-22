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
blog-generator
```

in this repository's root directory.

## Customization

### Configure the CLI

Set the following constants in `cli/cli.go`:

```go
// data source for the blog
const RepoURL string = "git@github.com:username/repo.git"
// folder to download the repo into
const TmpFolder string = "./tmp"
// output folder
const DestFolder string = "./www"
```

### Confugure the Generator

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
