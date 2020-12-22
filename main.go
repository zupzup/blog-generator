package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zupzup/blog-generator/cli"
)

func main() {
	serveFlag := flag.Bool("s", false, "start a local HTTP server to view static files for local testing")

	flag.Parse()

	argNums := len(os.Args)

	if argNums < 2 {
		cli.Run()
	} else {
		if *serveFlag {
			cli.Serve()
		} else {
			fmt.Println("Improper arguments")
		}
	}
}
