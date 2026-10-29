package main

import (
	"context"
	"io"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	urls := []string{
		"https://example.com",
		"https://example.org",
		"https://example.net",
	}

	results := make([]string, len(urls))

	for i, url := range urls {
		i, url := i, url
		g.Go(func() error {
			content, err := fetchURL(url)
			if err != nil {
				return err
			}
			results[i] = content
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatal("Fetch failed:", err)
	}

	_ = ctx
}

func fetchURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}
