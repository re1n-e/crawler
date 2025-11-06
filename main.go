package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	base_url := args[1]
	fmt.Printf("starting crawl of: %s\n", base_url)

	fmt.Println(getHTML(base_url))
}
