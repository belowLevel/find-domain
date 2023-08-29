package main

import (
	"bufio"
	"github.com/1121170088/find-domain/search"
	"os"
	"strings"
)

func main()  {
	search.Init("public_suffix_list.dat")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		domain := search.Search(line)
		if domain == "" {
			continue
		}
		os.Stdout.WriteString(domain)
		os.Stdout.WriteString("\n")
	}

	if scanner.Err() != nil {
		// Handle error.
	}
}
