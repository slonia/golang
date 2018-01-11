package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		if !strings.HasPrefix("http://", url) {
			url = "http://" + url
		}
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("fetch: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("status: %v\n", resp.Status)
		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Printf("fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}
