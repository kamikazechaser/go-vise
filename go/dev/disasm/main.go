package main

import (
	"fmt"
	"os"
	"io/ioutil"

	"git.defalsify.org/festive/vm"
)

func main() {
	if (len(os.Args) < 2) {
		os.Exit(1)
	}
	fp := os.Args[1]
	v, err := ioutil.ReadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read error: %v", err)
		os.Exit(1)
	}
	r, err := vm.ToString(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v", err)
		os.Exit(1)
	}
	fmt.Printf(r)
}
