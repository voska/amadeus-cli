package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	fmt.Fprintf(os.Stderr, "amadeus %s\n", version)
	os.Exit(0)
}
