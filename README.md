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
* File-Based Configuration

## Installation

```bash
go get github.com/zupzup/blog-generator
```

## Usage & Customization

### Configuration

The tool can be configured using a config file called `bloggen.yml`. There is a `bloggen.dist.yml` in the repository you can use as a template.

Example Config File:

```yml
generator:
    repo: 'https://github.com/zupzup/blog'
    tmp: 'tmp'
    dest: 'www'
    userss: true
blog:
    url: 'https://www.zupzup.org'
    language: 'en-us'
    description: 'A blog about Go, JavaScript, Open Source and Programming in General'
    dateformat: '02.01.2006'
    title: 'zupzup'
    author: 'Mario Zupan'
    frontpageposts: 10
    statics:
        files:
            - src: 'static/favicon.ico'
              dest: 'favicon.ico'
            - src: 'static/robots.txt'
              dest: 'robots.txt'
            - src: 'static/about.png'
              dest: 'about.png'
        templates:
            - src: 'static/about.html'
              dest: 'about'
```

### Running

Just execute

```bash
blog-generator
```

### Templates

Edit templates in `static` folder to your needs.

## Example Blog Repository

[Blog](https://github.com/zupzup/blog)
