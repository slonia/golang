package main

import (
	"bufio"
	"fmt"
	"os"
)

var words map[string][]string

func main() {
	counts := make(map[string]int)
	words = make(map[string][]string)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\t%v\n", n, line, words[line])
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		text := input.Text()
		counts[text]++
		_, ok := words[text]
		if !ok {
			words[text] = make([]string, 0)
		}
		words[text] = append(words[text], f.Name())
	}
}
