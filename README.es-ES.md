[![Go Report Card](https://goreportcard.com/badge/github.com/zupzup/calories)](https://goreportcard.com/report/github.com/zupzup/calories)

# blog-generator

Un generador estático de blogs que utiliza un repositorio de GitHub configurable como fuente de datos. Los posts se escriben en markdown con metadatos YAML adjuntos. [Este](https://github.com/zupzup/blog) es un repositorio de ejemplo para el blog en [https://zupzup.org/](https://zupzup.org/).

## Características

* Listado
* Generador de Sitemaps
* Fuente RSS
* Resaltado de Código
* Archivo
* Páginas Estáticas Configurables
* Etiquetas 
* Configuración basada en archivos

## Instalación

```bash
go get github.com/zupzup/blog-generator
```

## Uso y Personalización

### Configuración

La herramienta puede configurarse utilizando un archivo de configuración llamado `bloggen.yml`. Hay un archivo `bloggen.dist.yml` en el repositorio que puedes usar como plantilla.

Ejemplo de archivo de configuración:

```yml
generator:
    repo: 'https://git.sr.ht/~zupzup/blog'
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

### Ejecución

Solo ejecuta:

```bash
blog-generator
```

### Servir

Para ejecutar un servidor local que sirva el blog en http://127.0.0.1:8000, ejecuta:

```bash
blog-generator -s
```

### Plantillas

Edita las plantillas en la carpeta `static` según tus necesidades.

## Repositorio de Ejemplo del Blog

[Blog](https://git.sr.ht/~zupzup/blog)
