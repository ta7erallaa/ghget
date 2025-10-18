// Package downloader
package downloader

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ta7eralla/ghget/config"
	"github.com/ta7eralla/ghget/internal"
)

type Downloader struct {
	client internal.Client
}

func NewDownloader(client internal.Client) *Downloader {
	return &Downloader{
		client: client,
	}
}

func (d *Downloader) DownloadFromConfig(cfg *config.Config, filenames []string) error {
	urls := d.BuildURLs(cfg, filenames)

	filenames = d.parseFilenames(filenames)

	for i, url := range urls {
		if err := d.DownloadFile(url, filenames[i]); err != nil {
			return fmt.Errorf("failed download %s: %w", filenames[i], err)
		}
		fmt.Printf("Downloaded %s\n", filenames[i])
	}

	return nil
}

func (d *Downloader) DownloadFile(url, filename string) error {
	body, err := d.client.FetchFile(url)
	if err != nil {
		return err
	}
	defer body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, body); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (d *Downloader) BuildURLs(cfg *config.Config, filenames []string) []string {
	publicURL := "https://raw.githubusercontent.com/%s/%s/refs/heads/%s/%s"

	urls := make([]string, 0, len(filenames))
	for _, filename := range filenames {
		urls = append(urls, fmt.Sprintf(publicURL, cfg.Name, cfg.Repo, cfg.Branch, filename))
	}
	return urls
}

func (d *Downloader) parseFilenames(filenames []string) []string {
	var files []string

	for _, file := range filenames {
		if file == "" {
			continue
		}
		file = strings.TrimSuffix(file, "/")
		if idx := strings.LastIndex(file, "/"); idx != -1 {
			files = append(files, file[idx+1:])
		} else {
			files = append(files, file)
		}

	}
	return files
}
